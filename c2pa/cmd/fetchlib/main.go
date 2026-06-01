// fetchlib downloads the prebuilt c2pa-rs library bundle for the current
// GOOS/GOARCH from a GitHub Release of duggaraju/c2pa-go and extracts it into
// the location the cgo flags expect (../c2pa-rs/target/release by default), so
// consumers can build the package without a Rust toolchain.
//
// Typical use from a checkout of this repo:
//
//	go run ./c2pa/cmd/fetchlib                # default version, default dest
//	go run ./c2pa/cmd/fetchlib -dest /tmp/c2pa-libs
//	go run ./c2pa/cmd/fetchlib -link dynamic -env
//
// Or directly without cloning:
//
//	go run github.com/duggaraju/c2pa-go/c2pa/cmd/fetchlib@latest \
//	    -dest ./c2pa-rs/target/release
package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// LatestTag is the release tag fetched when -version is not provided. It is
// updated together with the matching c2pa-rs submodule pin and module tag.
const LatestTag = "latest"

// repoSlug is the GitHub repository hosting the release binaries.
const repoSlug = "duggaraju/c2pa-go"

type linkMode string

const (
	linkModeStatic  linkMode = "static"
	linkModeDynamic linkMode = "dynamic"
)

func main() {
	tag := flag.String("version", LatestTag, "release tag to fetch (e.g. c2pa/vX.Y.Z)")
	dest := flag.String("dest", "c2pa-rs/target/release", "destination directory (created if missing)")
	osName := flag.String("os", runtime.GOOS, "target operating system")
	arch := flag.String("arch", runtime.GOARCH, "target architecture")
	link := flag.String("link", string(linkModeStatic), "link mode for suggested CGO env vars: static or dynamic")
	envOnly := flag.Bool("env", false, "do not download; print suggested CGO env vars for an existing -dest and exit")
	flag.Parse()

	mode, err := parseLinkMode(*link)
	if err != nil {
		fail(err)
	}

	abs, err := filepath.Abs(*dest)
	if err != nil {
		fail(err)
	}

	if *envOnly {
		printEnv(abs, *osName, mode)
		return
	}

	asset := releaseAssetName(*osName, *arch)
	dlURL := releaseDownloadURL(repoSlug, *tag, asset)

	if err := os.MkdirAll(abs, 0o755); err != nil {
		fail(err)
	}

	fmt.Fprintf(os.Stderr, "fetchlib: downloading %s\n", dlURL)
	if err := downloadAndExtract(dlURL, abs); err != nil {
		fail(err)
	}

	fmt.Fprintf(os.Stderr, "fetchlib: extracted to %s\n\n", abs)
	fmt.Fprintf(os.Stderr, "To build against the prebuilt library, either:\n")
	fmt.Fprintf(os.Stderr, "  1) Use the default cgo flags: ensure -dest matches\n")
	fmt.Fprintf(os.Stderr, "     <repo>/c2pa-rs/target/release relative to the c2pa\n")
	fmt.Fprintf(os.Stderr, "     package source, then `go build -tags=release ./...`\n")
	fmt.Fprintf(os.Stderr, "  2) Or set CGO_*FLAGS yourself with -link=%s, e.g.:\n\n", mode)
	printEnv(abs, *osName, mode)
}

func parseLinkMode(value string) (linkMode, error) {
	switch linkMode(strings.ToLower(strings.TrimSpace(value))) {
	case linkModeStatic:
		return linkModeStatic, nil
	case linkModeDynamic:
		return linkModeDynamic, nil
	default:
		return "", fmt.Errorf("unsupported -link value %q (want static or dynamic)", value)
	}
}

func releaseAssetName(osName, arch string) string {
	return fmt.Sprintf("c2pa-c-libs-release-%s-%s.tar.gz", osName, arch)
}

func releaseDownloadURL(repoSlug, tag, asset string) string {
	tag = strings.TrimSpace(tag)
	if strings.EqualFold(tag, LatestTag) {
		return fmt.Sprintf("https://github.com/%s/releases/latest/download/%s",
			repoSlug, asset)
	}

	return fmt.Sprintf("https://github.com/%s/releases/download/%s/%s",
		repoSlug, url.PathEscape(tag), asset)
}

func downloadAndExtract(dlURL, dest string) error {
	resp, err := http.Get(dlURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Sanitize entry path to prevent directory traversal (zip-slip).
		name := filepath.Clean(hdr.Name)
		if name == "." || name == "/" {
			continue
		}
		if filepath.IsAbs(name) || strings.HasPrefix(name, "..") ||
			strings.Contains(name, string(filepath.Separator)+"..") {
			return fmt.Errorf("unsafe entry in archive: %s", hdr.Name)
		}
		out := filepath.Join(dest, name)

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(out, 0o755); err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
				return err
			}
			f, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			if err := f.Close(); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "  extracted %s\n", name)
		default:
			// Skip symlinks and other entry types; the bundle should not need them.
		}
	}
	return nil
}

func printEnv(dest, osName string, mode linkMode) {
	for _, line := range envLines(dest, osName, mode) {
		fmt.Println(line)
	}
}

func envLines(dest, osName string, mode linkMode) []string {
	lines := []string{fmt.Sprintf("export CGO_CFLAGS=\"-I%s\"", dest)}

	switch osName {
	case "windows":
		// MSYS2 / mingw style; users may need to translate paths for cmd.exe.
		if mode == linkModeDynamic {
			lines = append(lines,
				fmt.Sprintf("export CGO_LDFLAGS=\"-L%s -l:c2pa_c.lib -lws2_32 -luserenv -ladvapi32 -lncrypt -lcrypt32 -lbcrypt -lsecur32 -lntdll -lkernel32 -lole32 -loleaut32 -lpsapi -liphlpapi\"", dest),
				fmt.Sprintf("export PATH=\"%s:$PATH\"", dest),
			)
			return lines
		}

		return append(lines,
			fmt.Sprintf("export CGO_LDFLAGS=\"-L%s -l:libc2pa_c.a -lws2_32 -luserenv -ladvapi32 -lncrypt -lcrypt32 -lbcrypt -lsecur32 -lntdll -lkernel32 -lole32 -loleaut32 -lpsapi -liphlpapi\"", dest),
		)
	case "darwin":
		if mode == linkModeDynamic {
			return append(lines,
				fmt.Sprintf("export CGO_LDFLAGS=\"-L%s -Wl,-rpath,%s -lc2pa_c -framework Security -framework CoreFoundation -framework SystemConfiguration -lresolv -ldl -lm\"", dest, dest),
			)
		}

		return append(lines,
			fmt.Sprintf("export CGO_LDFLAGS=\"%s/libc2pa_c.a -framework Security -framework CoreFoundation -framework SystemConfiguration -lresolv -ldl -lm\"", dest),
		)
	default:
		if mode == linkModeDynamic {
			return append(lines,
				fmt.Sprintf("export CGO_LDFLAGS=\"-L%s -Wl,-rpath,%s -Wl,-Bdynamic -lc2pa_c -lm -ldl -lpthread\"", dest, dest),
			)
		}

		return append(lines,
			fmt.Sprintf("export CGO_LDFLAGS=\"-L%s -Wl,-Bstatic -lc2pa_c -Wl,-Bdynamic -lm -ldl -lpthread\"", dest),
		)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "fetchlib:", err)
	os.Exit(1)
}
