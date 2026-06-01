package main

import (
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
