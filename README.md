# c2pa-go

[![Go Reference](https://pkg.go.dev/badge/github.com/duggaraju/c2pa-go/c2pa.svg)](https://pkg.go.dev/github.com/duggaraju/c2pa-go/c2pa)
[![Go Report Card](https://goreportcard.com/badge/github.com/duggaraju/c2pa-go/c2pa)](https://goreportcard.com/report/github.com/duggaraju/c2pa-go/c2pa)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![c2pa-rs](https://img.shields.io/badge/c2pa--rs-v0.83.0-orange?logo=rust)](https://github.com/contentauth/c2pa-rs)

Go bindings for the [c2pa-rs](https://github.com/contentauth/c2pa-rs) Content
Authenticity SDK. This package exposes the C2PA reader, builder, signer and
settings APIs through cgo so that Go programs can read, create and sign
[C2PA](https://c2pa.org) manifests (also known as Content Credentials), adding
Content Authenticity and Content Integrity metadata to media assets.

**Keywords:** c2pa, c2pa-go,  c2pa-rs, Content Authenticity, Content Integrity, Content Credentials, content provenance.

The bindings track c2pa-rs `0.83.x` and link against the `c2pa_c` shared library
produced by the `c2pa_c_ffi` crate.

## Status

- Supported assets: anything supported by c2pa-rs (JPEG, PNG, TIFF, MP4, etc.).
- Supported algorithms: `Ps256`/`Ps384`/`Ps512`, `Es256`/`Es384`/`Es512`,
  `Ed25519`.
- Supported platforms (CI-built):
  - Linux amd64, Linux arm64
  - macOS amd64 (Intel), macOS arm64 (Apple Silicon)
  - Windows amd64, Windows arm64 (mingw / clang-mingw toolchain via cgo)
- Bindings are not yet API-stable; the surface may change before a tagged
  release.

## Repository layout

```
c2pa/           Go package (module github.com/duggaraju/c2pa-go/c2pa)
c2pa/schema/    Generated Go types from the c2pa JSON schemas
                (package github.com/duggaraju/c2pa-go/c2pa/schema)
example/        Sample CLI that reads and signs assets
c2pa-rs/        Pinned upstream submodule built into the c2pa_c static lib
```

## Requirements

- Go 1.24+
- `cgo` enabled (`CGO_ENABLED=1`) with a C toolchain:
  - Linux: `gcc` or `clang`
  - macOS: Xcode Command Line Tools (`xcode-select --install`)
  - Windows: a mingw-w64 toolchain on `PATH` (e.g. MSYS2 `MINGW64` for amd64,
    `CLANGARM64` for arm64). The MSVC `cl.exe` toolchain is **not** supported
    by cgo.
- Rust toolchain (only required to rebuild `c2pa-rs`)
  - On Windows, install the **gnu** variant (`stable-x86_64-pc-windows-gnu` on
    amd64 or `stable-aarch64-pc-windows-gnullvm` on arm64) so `cargo` emits a
    mingw-compatible `libc2pa_c.a`. Setting it as your default with
    `rustup default stable-<triple>` keeps `cargo build` writing to
    `c2pa-rs/target/{debug,release}/`.
- `git submodule update --init --recursive` for the pinned c2pa-rs checkout
- `node`/`npx` (only required to regenerate the `c2pa/schema/` package)

## Building

The repo uses `go:generate` for every build artifact that isn't pure Go. The
cgo flag files in [c2pa](c2pa) target the host OS automatically via Go build
constraints; cross-arch builds are not currently supported by the
`go:generate` directives — build on a host matching your target `GOOS/GOARCH`
or download a prebuilt library from a GitHub Release (see below).

```sh
git submodule update --init --recursive

# Build the c2pa-rs C FFI static library (debug flavor) into c2pa-rs/target/debug.
go generate ./c2pa/...

# Or build the release flavor into c2pa-rs/target/release.
go generate -tags=release ./c2pa/...

# Build the bindings and the example CLI.
go build ./c2pa                       # debug
go build -tags=release ./c2pa         # release
(cd example && go build .)
```

The `go:generate` directives in [c2pa/generate.go](c2pa/generate.go),
[c2pa/generate_debug.go](c2pa/generate_debug.go) and
[c2pa/generate_release.go](c2pa/generate_release.go) regenerate the
[c2pa/schema](c2pa/schema) package and run `cargo build --features rust_native_crypto,http`
inside `c2pa-rs/c2pa_c_ffi` for the flavor matching the `-tags` passed to
`go generate`. The cgo flag files in [c2pa](c2pa) point the linker at
`c2pa-rs/target/{debug,release}` accordingly.

### Using prebuilt C libraries

CI publishes `libc2pa_c` + `c2pa.h` for every supported OS/arch as GitHub
Release assets (named `c2pa-c-libs-release-<os>-<arch>-<sha>.tar.gz`). To
consume the bindings without a Rust toolchain:

1. Download the archive matching your `GOOS`/`GOARCH` from the
   [Releases page](https://github.com/duggaraju/c2pa-go/releases).
2. Extract its contents into `c2pa-rs/target/release/` (or `target/debug/` for
   the debug flavor) inside your checkout of this repo.
3. `go build -tags=release ./c2pa` will then link against the prebuilt
   library; you no longer need `cargo` or `rustup` installed.

### Regenerating the schema package

The Go types in [c2pa/schema/schema.go](c2pa/schema/schema.go) are produced from
the JSON schemas emitted by the upstream `export_schema` Rust binary. The
generated file is checked in so that consumers of this module do not need Rust
or Node installed. To refresh it after bumping c2pa-rs:

```sh
go generate ./c2pa/schema/...
```

This runs `cargo run` for `c2pa-rs/export_schema`, rewrites the top-level
schema titles to `ManifestDefinition` and `ManifestStore`, and runs
`quicktype` (via `npx`) to emit [c2pa/schema/schema.go](c2pa/schema/schema.go).

In your own Go module, add the package as a dependency:

```sh
go get github.com/duggaraju/c2pa-go/c2pa
```

You will still need a built `libc2pa_c` available at link time. The cgo flag
files in [c2pa](c2pa) point at `c2pa-rs/target/{debug,release}`; adjust them for
your deployment if you ship the shared library elsewhere.

## Usage

All operations are scoped to a `*c2pa.Context`. A context optionally carries
`Settings` (trust list, verification policy, builder defaults, etc.).

### Reading a manifest

```go
package main

import (
    "fmt"
    "log"

    c2pa "github.com/duggaraju/c2pa-go/c2pa"
)

func main() {
    ctx, err := c2pa.NewContext()
    if err != nil {
        log.Fatal(err)
    }
    defer ctx.Close()

    r, err := c2pa.ReaderFromFile(ctx, "signed.jpg")
    if err != nil {
        log.Fatal(err)
    }
    defer r.Close()

    fmt.Println(r.Json())
    fmt.Println("embedded:", r.IsEmbedded())
}
```

### Loading settings from TOML

```go
builder, err := c2pa.NewContextBuilder()
if err != nil { log.Fatal(err) }
defer builder.Close()

content, _ := os.ReadFile("settings.toml")
settings, err := c2pa.NewSettings()
if err != nil { log.Fatal(err) }
defer settings.Close()

if err := settings.UpdateFromString(string(content), "toml"); err != nil {
    log.Fatal(err)
}
if err := builder.SetSettings(settings); err != nil {
    log.Fatal(err)
}

ctx, err := builder.Build()
if err != nil { log.Fatal(err) }
defer ctx.Close()
```

### Signing an asset

Provide a manifest definition and an implementation of `c2pa.Signer`. The
binding will invoke `Sign` via a cgo callback when c2pa-rs needs a COSE
signature over the claim bytes.

```go
manifestJson := `{
  "claim_generator": "my-app/1.0",
  "title": "example.jpg"
}`

builder, err := c2pa.BuilderFromJson(ctx, manifestJson)
if err != nil { log.Fatal(err) }
defer builder.Close()

signer := &MyPs256Signer{ /* certs PEM + RSA key */ }

manifest, err := builder.Sign("in.jpg", "out.jpg", signer)
if err != nil { log.Fatal(err) }
fmt.Printf("wrote %d manifest bytes\n", len(manifest))
```

A minimal `Signer` implementation (PS256, RSA-PSS over SHA-256):

```go
type MyPs256Signer struct {
    certsPEM string
    key      *rsa.PrivateKey
}

func (s *MyPs256Signer) Alg() c2pa.SigningAlg     { return c2pa.SigningAlgPs256 }
func (s *MyPs256Signer) Certificates() string    { return s.certsPEM }
func (s *MyPs256Signer) TimeStampUrl() string    { return "" }

func (s *MyPs256Signer) Sign(input, output []byte) (int, error) {
    h := sha256.Sum256(input)
    sig, err := rsa.SignPSS(rand.Reader, s.key, crypto.SHA256, h[:],
        &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash})
    if err != nil {
        return -1, err
    }
    if len(sig) > len(output) {
        return -1, fmt.Errorf("output too small")
    }
    return copy(output, sig), nil
}
```

The certificate PEM bundle (leaf first, then intermediates) is what
c2pa-rs embeds in the COSE `x5chain` header.

### Builder helpers

Beyond `Sign`, `Builder` exposes most of the upstream API:

```go
builder.SetIntent(c2pa.IntentCreate, c2pa.SourceDigitalCapture)
builder.SetRemoteUrl("https://example.com/manifest.c2pa")
builder.SetNoEmbed()
builder.AddAction(`{"action":"c2pa.color_adjustments"}`)
builder.AddResourceFromFile("thumbnail", "thumb.jpg")
builder.AddIngredientFromFile(`{"title":"parent"}`, "parent.jpg")
builder.ToArchiveFile("manifest.c2pa")
```

### Reader helpers

```go
fmt.Println(r.Json())          // summary JSON
fmt.Println(r.DetailedJson())  // detailed JSON
fmt.Println(r.RemoteUrl())     // "" if remote-only
r.ResourceToFile("self#jumbf=c2pa.assertions/c2pa.thumbnail", "thumb.jpg")
```

### Typed APIs

In addition to the raw JSON entry points, the bindings expose Go types
generated from the upstream JSON schemas in the
[`schema`](https://pkg.go.dev/github.com/duggaraju/c2pa-go/c2pa/schema)
subpackage. Use them when you would rather build manifests, parse manifest
stores, or configure settings with statically typed values instead of hand-
written JSON strings.

```go
import (
    c2pa "github.com/duggaraju/c2pa-go/c2pa"
    "github.com/duggaraju/c2pa-go/c2pa/schema"
)
```

#### Typed manifest definition

`BuilderFromDefinition` accepts a `*schema.ManifestDefinition` and marshals it
internally before handing it to c2pa-rs:

```go
title := "example.jpg"
format := "image/jpeg"
vendor := "my-app"

def := &schema.ManifestDefinition{
    Title:  &title,
    Format: &format,
    Vendor: &vendor,
    ClaimGeneratorInfo: []schema.ClaimGeneratorInfoElement{{
        Name:    "my-app",
        Version: ptr("1.0.0"),
    }},
    Assertions: []schema.AssertionElement{{
        Label: "c2pa.actions",
        Data: map[string]any{
            "actions": []map[string]any{{"action": "c2pa.created"}},
        },
    }},
}

builder, err := c2pa.BuilderFromDefinition(ctx, def)
if err != nil { log.Fatal(err) }
defer builder.Close()

// Add an extra action assertion from a typed value.
builder.AddActionTyped(map[string]any{"action": "c2pa.color_adjustments"})
```

#### Typed manifest store (Reader)

`Reader.Manifest` and `Reader.DetailedManifest` return a parsed
`*schema.ManifestStore` rather than a JSON string:

```go
store, err := r.Manifest()
if err != nil { log.Fatal(err) }

if store.ActiveManifest != nil {
    active := store.Manifests[*store.ActiveManifest]
    if active.Title != nil {
        fmt.Println("title:", *active.Title)
    }
    for _, a := range active.Assertions {
        fmt.Println("assertion:", a.Label)
    }
}

if store.ValidationState != nil {
    fmt.Println("state:", *store.ValidationState)
}
```

#### Typed settings

`Settings.UpdateFrom` accepts a `*schema.Settings` value:

```go
settings, err := c2pa.NewSettings()
if err != nil { log.Fatal(err) }
defer settings.Close()

typed := &schema.Settings{
    Verify: &schema.Verify{
        VerifyAfterReading: ptr(true),
        VerifyTrust:        ptr(true),
    },
}
if err := settings.UpdateFrom(typed); err != nil {
    log.Fatal(err)
}
```

A tiny helper for taking the address of a literal is convenient since most
schema fields are pointers (so the JSON `omitempty` semantics are preserved):

```go
func ptr[T any](v T) *T { return &v }
```

## Example CLI

The [example](example) directory contains a small CLI demonstrating `read` and
`sign`:

```sh
cd example
go build .

./example -v
./example read -i ../c2pa-rs/sdk/tests/fixtures/C.jpg
./example sign \
    -i ../c2pa-rs/sdk/tests/fixtures/C.jpg \
    -o /tmp/signed.jpg \
    -c ../c2pa-rs/sdk/tests/fixtures/certs/ps256.pub \
    -k ../c2pa-rs/sdk/tests/fixtures/certs/ps256.pem
```

Both subcommands accept `-s settings.toml` to load a settings file; see
[example/settings.toml](example/settings.toml) for a working sample with a
custom trust list and verification policy.

## License

This project is dual-licensed under the MIT license and the Apache License 2.0
(matching c2pa-rs). See [LICENSE](LICENSE) for details.
