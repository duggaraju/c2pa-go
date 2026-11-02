package main

import (
	"runtime"
	"strings"
	"testing"
)

func TestParseLinkMode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    linkMode
		wantErr bool
	}{
		{name: "default static", input: "static", want: linkModeStatic},
		{name: "case insensitive dynamic", input: "Dynamic", want: linkModeDynamic},
		{name: "reject invalid", input: "shared", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseLinkMode(test.input)
			if test.wantErr {
				if err == nil {
					t.Fatalf("parseLinkMode(%q) error = nil, want error", test.input)
				}
				return
			}

			if err != nil {
				t.Fatalf("parseLinkMode(%q) error = %v", test.input, err)
			}
			if got != test.want {
				t.Fatalf("parseLinkMode(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}

func TestEnvLinesLinux(t *testing.T) {
	tests := []struct {
		name      string
		mode      linkMode
		wantParts []string
		notParts  []string
	}{
		{
			name:      "static",
			mode:      linkModeStatic,
			wantParts: []string{"-Wl,-Bstatic -lc2pa_c -Wl,-Bdynamic", "export CGO_CFLAGS=\"-I/tmp/c2pa\""},
			notParts:  []string{"-Wl,-rpath,/tmp/c2pa"},
		},
		{
			name:      "dynamic",
			mode:      linkModeDynamic,
			wantParts: []string{"-Wl,-rpath,/tmp/c2pa", "-Wl,-Bdynamic -lc2pa_c"},
			notParts:  []string{"-Wl,-Bstatic -lc2pa_c"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := strings.Join(envLines("/tmp/c2pa", "linux", test.mode), "\n")
			for _, part := range test.wantParts {
				if !strings.Contains(got, part) {
					t.Fatalf("envLines() missing %q in %q", part, got)
				}
			}
			for _, part := range test.notParts {
				if strings.Contains(got, part) {
					t.Fatalf("envLines() unexpectedly contains %q in %q", part, got)
				}
			}
		})
	}
}

func TestReleaseDownloadURL(t *testing.T) {
	tests := []struct {
		name string
		tag  string
		want string
	}{
		{
			name: "latest sentinel",
			tag:  "latest",
			want: "https://github.com/duggaraju/c2pa-go/releases/latest/download/c2pa-c-libs-release-linux-amd64.tar.gz",
		},
		{
			name: "explicit tag",
			tag:  "c2pa/v0.85.1",
			want: "https://github.com/duggaraju/c2pa-go/releases/download/c2pa%2Fv0.85.1/c2pa-c-libs-release-linux-amd64.tar.gz",
		},
		{
			name: "latest case insensitive",
			tag:  " Latest ",
			want: "https://github.com/duggaraju/c2pa-go/releases/latest/download/c2pa-c-libs-release-linux-amd64.tar.gz",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := releaseDownloadURL("duggaraju/c2pa-go", test.tag, "c2pa-c-libs-release-linux-amd64.tar.gz")
			if got != test.want {
				t.Fatalf("releaseDownloadURL() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestCheckLibc(t *testing.T) {
	musl := func() bool { return true }
	glibc := func() bool { return false }

	tests := []struct {
		name    string
		libc    string
		os      string
		arch    string
		detect  func() bool
		wantErr bool
	}{
		{name: "auto on glibc host", libc: "auto", os: runtime.GOOS, arch: runtime.GOARCH, detect: glibc},
		{name: "auto on musl host", libc: "auto", os: runtime.GOOS, arch: runtime.GOARCH, detect: musl, wantErr: runtime.GOOS == "linux"},
		{name: "auto ignores host when cross targeting", libc: "auto", os: "linux", arch: "not-" + runtime.GOARCH, detect: musl},
		{name: "explicit musl on linux", libc: "musl", os: "linux", arch: "amd64", detect: glibc, wantErr: true},
		{name: "explicit glibc overrides detection", libc: "glibc", os: "linux", arch: runtime.GOARCH, detect: musl},
		{name: "musl is irrelevant off linux", libc: "musl", os: "darwin", arch: "arm64", detect: musl},
		{name: "case insensitive", libc: " MUSL ", os: "linux", arch: "amd64", detect: glibc, wantErr: true},
		{name: "reject invalid", libc: "uclibc", os: "linux", arch: "amd64", detect: glibc, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := checkLibc(test.libc, test.os, test.arch, test.detect)
			if (err != nil) != test.wantErr {
				t.Fatalf("checkLibc(%q, %q, %q) error = %v, wantErr %v", test.libc, test.os, test.arch, err, test.wantErr)
			}
		})
	}
}

func TestLddOutputIsMusl(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   bool
	}{
		{name: "musl", output: "musl libc (x86_64)\nVersion 1.2.5\nDynamic Program Loader\n", want: true},
		{name: "glibc", output: "ldd (Ubuntu GLIBC 2.43-2ubuntu2.4) 2.43\n", want: false},
		{name: "missing ldd", output: "", want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := lddOutputIsMusl(test.output); got != test.want {
				t.Fatalf("lddOutputIsMusl(%q) = %v, want %v", test.output, got, test.want)
			}
		})
	}
}
