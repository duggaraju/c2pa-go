// fetchlib downloads the prebuilt c2pa-rs static library bundle for the
// current GOOS/GOARCH from a GitHub Release of duggaraju/c2pa-go and extracts
// it into the location the cgo flags expect (../c2pa-rs/target/release by
// default), so consumers can build the package without a Rust toolchain.
//
// Typical use from a checkout of this repo:
//
//	go run ./c2pa/cmd/fetchlib                # default version, default dest
//	go run ./c2pa/cmd/fetchlib -version c2pa/v0.84.1
//	go run ./c2pa/cmd/fetchlib -dest /tmp/c2pa-libs
//
// Or directly without cloning:
//
//	go run github.com/duggaraju/c2pa-go/c2pa/cmd/fetchlib@v0.84.1 \
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

// DefaultTag is the release tag fetched when -version is not provided. It is
// updated together with the matching c2pa-rs submodule pin and module tag.
const DefaultTag = "c2pa/v0.84.1"

// repoSlug is the GitHub repository hosting the release binaries.
const repoSlug = "duggaraju/c2pa-go"

func main() {
	tag := flag.String("version", DefaultTag, "release tag to fetch (e.g. c2pa/v0.84.1)")
	dest := flag.String("dest", "c2pa-rs/target/release", "destination directory (created if missing)")
	osName := flag.String("os", runtime.GOOS, "target operating system")
	arch := flag.String("arch", runtime.GOARCH, "target architecture")
	envOnly := flag.Bool("env", false, "do not download; print suggested CGO env vars for an existing -dest and exit")
	flag.Parse()

	abs, err := filepath.Abs(*dest)
	if err != nil {
		fail(err)
	}

	if *envOnly {
		printEnv(abs, *osName)
		return
	}

	asset := fmt.Sprintf("c2pa-c-libs-release-%s-%s.tar.gz", *osName, *arch)
	dlURL := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s",
		repoSlug, url.PathEscape(*tag), asset)

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
	fmt.Fprintf(os.Stderr, "  2) Or set CGO_*FLAGS yourself, e.g.:\n\n")
	printEnv(abs, *osName)
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

func printEnv(dest, osName string) {
	switch osName {
	case "windows":
		// MSYS2 / mingw style; users may need to translate paths for cmd.exe.
		fmt.Printf("export CGO_CFLAGS=\"-I%s\"\n", dest)
		fmt.Printf("export CGO_LDFLAGS=\"-L%s -lc2pa_c -lws2_32 -lbcrypt -luserenv -lntdll\"\n", dest)
	case "darwin":
		fmt.Printf("export CGO_CFLAGS=\"-I%s\"\n", dest)
		fmt.Printf("export CGO_LDFLAGS=\"-L%s -lc2pa_c -framework Security -framework CoreFoundation\"\n", dest)
	default:
		fmt.Printf("export CGO_CFLAGS=\"-I%s\"\n", dest)
		fmt.Printf("export CGO_LDFLAGS=\"-L%s -lc2pa_c -lm -ldl -lpthread\"\n", dest)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "fetchlib:", err)
	os.Exit(1)
}
