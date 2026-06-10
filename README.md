# c2pa-go

[![Go Reference](https://pkg.go.dev/badge/github.com/duggaraju/c2pa-go/c2pa.svg)](https://pkg.go.dev/github.com/duggaraju/c2pa-go/c2pa)
[![Go Report Card](https://goreportcard.com/badge/github.com/duggaraju/c2pa-go/c2pa)](https://goreportcard.com/report/github.com/duggaraju/c2pa-go/c2pa)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![c2pa-rs](https://img.shields.io/badge/c2pa--rs-v0.86.1-orange?logo=rust)](https://github.com/contentauth/c2pa-rs)

Go bindings for the [c2pa-rs](https://github.com/contentauth/c2pa-rs) Content
Authenticity SDK. This package exposes the C2PA reader, builder, signer and
settings APIs through cgo so that Go programs can read, create and sign
[C2PA](https://c2pa.org) manifests (also known as Content Credentials), adding
Content Authenticity and Content Integrity metadata to media assets.

**Keywords:** c2pa, c2pa-go,  c2pa-rs, Content Authenticity, Content Integrity, Content Credentials, content provenance.

The bindings track c2pa-rs `0.86.x` and link against the `c2pa_c` shared library
produced by the `c2pa_c_ffi` crate.

## Status

- Supported assets: anything supported by c2pa-rs (JPEG, PNG, TIFF, MP4, etc.).
- Supported algorithms: `Ps256`/`Ps384`/`Ps512`, `Es256`/`Es384`/`Es512`,
  `Ed25519`.
- Supported platforms (CI-built):
  - Linux amd64, Linux arm64
  - macOS amd64 (Intel), macOS arm64 (Apple Silicon)
  - Windows amd64 (mingw / gcc toolchain via cgo)
  - Windows arm64: temporarily disabled in CI while the
    `aarch64-pc-windows-gnullvm` toolchain story is sorted out.
- Bindings are not yet API-stable; the surface may change before a tagged
  release.

## Repository layout

```
c2pa/                 Go package (module github.com/duggaraju/c2pa-go/c2pa)
c2pa/schema/          Generated Go types from the c2pa JSON schemas
                      (package github.com/duggaraju/c2pa-go/c2pa/schema)
c2pa/cmd/fetchlib/    Helper CLI that downloads the prebuilt native library
                      tarball for the current OS/arch from a GitHub Release
                      and extracts it into c2pa-rs/target/release (see
                      "Using prebuilt C libraries" below)
c2pa-cli/             Sample CLI that reads and signs assets
                      (module github.com/duggaraju/c2pa-go/c2pa-cli)
c2pa-rs/              Pinned upstream submodule built into the c2pa_c static lib
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

The `c2pa` package also compiles cleanly with `CGO_ENABLED=0` — every cgo
call is isolated in [c2pa/native.go](c2pa/native.go) and stubbed in
[c2pa/native_nocgo.go](c2pa/native_nocgo.go). This is purely so that
pkg.go.dev and other doc/indexing tools can build the package; at runtime
every constructor returns a "native library unavailable" error.

## Building

The repo uses `go:generate` for every build artifact that isn't pure Go. The
cgo flag files in [c2pa](c2pa) target the host OS automatically via Go build
constraints; cross-arch builds are not currently supported by the
`go:generate` directives — build on a host matching your target `GOOS/GOARCH`
or download a prebuilt library from a GitHub Release (see below).

There are two supported paths:

1. **Without the Rust toolchain** (recommended for consumers) — download a
   prebuilt `libc2pa_c` from a GitHub Release with the `fetchlib` helper, then
   run `go build`. No `cargo`/`rustup` required. See
   [Using prebuilt C libraries](#using-prebuilt-c-libraries) below for the
   full recipe; the short version is:

   ```sh
   git submodule update --init --recursive
   go run ./c2pa/cmd/fetchlib            # populates c2pa-rs/target/release/
   go build -tags=release ./c2pa
   (cd c2pa-cli && go build -tags=release .)
   ```

   From a project that just depends on this module (no checkout needed):

   ```sh
    go run github.com/duggaraju/c2pa-go/c2pa/cmd/fetchlib@latest \
       -dest ./c2pa-rs/target/release
   go build -tags=release ./...
   ```

2. **With the Rust toolchain** — rebuild the C FFI from source via
   `go generate`:

   ```sh
   git submodule update --init --recursive

   # Build the c2pa-rs C FFI static library (debug flavor) into c2pa-rs/target/debug.
   go generate ./c2pa/...

   # Or build the release flavor into c2pa-rs/target/release.
   go generate -tags=release ./c2pa/...

   # Build the bindings and the sample CLI.
   go build ./c2pa                       # debug
   go build -tags=release ./c2pa         # release
   (cd c2pa-cli && go build .)
   ```

The `go:generate` directives in [c2pa/generate.go](c2pa/generate.go),
[c2pa/generate_debug.go](c2pa/generate_debug.go) and
[c2pa/generate_release.go](c2pa/generate_release.go) regenerate the
[c2pa/schema](c2pa/schema) package and run `cargo build --features rust_native_crypto`
inside `c2pa-rs/c2pa_c_ffi` for the flavor matching the `-tags` passed to
`go generate`. The cgo flag files in [c2pa](c2pa) point the linker at
`c2pa-rs/target/{debug,release}` accordingly.

Remote operations (remote manifest fetches, OCSP, RFC 3161 timestamp
authorities) are handled by the bundled reqwest-based HTTP client in
`libc2pa_c` and work out of the box. They can optionally be redirected to
a Go-side resolver wired up through `ContextBuilder.SetHttpResolver` — see
[HTTP resolver](#http-resolver) below — if you want to control transport,
proxy through a custom `http.Client`, inject test fixtures, or eventually
build `libc2pa_c` without the `http` feature.

### Using prebuilt C libraries

CI publishes `libc2pa_c` + `c2pa.h` for every supported OS/arch as GitHub
Release assets (named `c2pa-c-libs-release-<os>-<arch>.tar.gz`). Each bundle
contains the static archive plus the platform's shared library/import library
when available. To consume the bindings without a Rust toolchain, use the
`fetchlib` helper to download and extract the bundle for your host:

```sh
# From a checkout of this repo:
go run ./c2pa/cmd/fetchlib                              # default version, default dest
go run ./c2pa/cmd/fetchlib -version c2pa/vX.Y.Z         # pin a specific release
go run ./c2pa/cmd/fetchlib -link dynamic -env           # print dynamic-linker CGO_*FLAGS
go run ./c2pa/cmd/fetchlib -dest /tmp/c2pa-libs -env    # also print suggested CGO_*FLAGS

# Or from an arbitrary working directory, without cloning the repo:
go run github.com/duggaraju/c2pa-go/c2pa/cmd/fetchlib@latest \
    -dest ./c2pa-rs/target/release
```

By default the helper writes to `c2pa-rs/target/release/` relative to the
current directory, which matches the cgo flag files. `fetchlib` also accepts
`-link static|dynamic`; `static` is the default and matches the Linux release
cgo flags in this repo. After running it, `go build -tags=release ./c2pa`
will link against the prebuilt library and you no longer need `cargo` or
`rustup` installed.

If you ship the library to a non-default location, pass `-dest` and either
update the cgo flag files in [c2pa](c2pa) or override `CGO_CFLAGS`/`CGO_LDFLAGS`
at build time (the helper prints a starter set with `-env`).

You can also download the archive manually from the
[Releases page](https://github.com/duggaraju/c2pa-go/releases) and extract it
into the same destination directory if you prefer to avoid running the helper.

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

    r, err := c2pa.NewReader(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer r.Close()

    if err := r.WithFile("signed.jpg"); err != nil {
        log.Fatal(err)
    }

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

Provide a manifest definition and a signer accepted by `Builder.Sign`.
For Go callback-based signing, implement `c2pa.CallbackSigner`; the binding
will invoke `Sign` via a cgo callback when c2pa-rs needs a COSE signature over
the claim bytes.

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

A minimal `CallbackSigner` implementation (PS256, RSA-PSS over SHA-256):

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

### HTTP resolver

The `c2pa_c` library is built with the reqwest-based `http` feature, so
remote-manifest fetches, OCSP lookups, and RFC 3161 timestamping (verified
against `http://timestamp.digicert.com`) all work without any extra wiring.

The bindings also expose an optional Go-side resolver. Use it when you want
to route every C2PA HTTP request through your own `http.Client` (custom
transport, proxy, request signing, fixtures in tests, etc.) or when you
plan to build `libc2pa_c` without the `http` feature to shrink the binary.
Install one on the `ContextBuilder` before calling `Build`:

```go
builder, err := c2pa.NewContextBuilder()
if err != nil { log.Fatal(err) }
defer builder.Close()

resolver := &c2pa.DefaultHttpResolver{
    // Optional; defaults to http.DefaultClient.
    Client: &http.Client{Timeout: 30 * time.Second},
}

if err := builder.SetHttpResolver(resolver); err != nil {
    log.Fatal(err)
}

ctx, err := builder.Build()
if err != nil { log.Fatal(err) }
defer ctx.Close()
```

For custom transport behaviour (proxying, request signing, fixtures in tests,
etc.), implement the `c2pa.HttpResolver` interface directly:

```go
type HttpResolver interface {
    Resolve(req *http.Request) (*http.Response, error)
}
```

The resolver is invoked synchronously from c2pa-rs; return errors normally
and they are surfaced through the usual `C2paError()` channel. Installing
a resolver fully replaces the bundled reqwest client for the lifetime of
the `Context`.

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

## Sample CLI

The [c2pa-cli](c2pa-cli) directory contains a small CLI demonstrating `read`
and `sign`. It is published as its own module,
`github.com/duggaraju/c2pa-go/c2pa-cli`, so it can be installed directly:

```sh
go install github.com/duggaraju/c2pa-go/c2pa-cli@latest
```

`go install` only links if a prebuilt `libc2pa_c` is reachable by the linker.
The easiest way to set that up — without a Rust toolchain or even a checkout
of this repo — is to use the `fetchlib` helper and the CGO env vars it prints:

```sh
# 1. Download the prebuilt c2pa_c library bundle to a directory of your choice.
go run github.com/duggaraju/c2pa-go/c2pa/cmd/fetchlib@latest -dest "$PWD/libs"

# 2. Export the CGO flags fetchlib prints. -env prints them without re-downloading.
eval "$(go run github.com/duggaraju/c2pa-go/c2pa/cmd/fetchlib@latest \
    -dest "$PWD/libs" -env)"

# 3. Install the CLI.
go install github.com/duggaraju/c2pa-go/c2pa-cli@latest
c2pa-cli -v
```

From a checkout of this repo the default cgo flags already point at
`c2pa-rs/target/release`, so the build is a two-step recipe:

```sh
go run ./c2pa/cmd/fetchlib            # populates c2pa-rs/target/release/
cd c2pa-cli
go build -tags=release .

./c2pa-cli -v
./c2pa-cli read -i ../c2pa-rs/sdk/tests/fixtures/C.jpg
./c2pa-cli sign \
    -i ../c2pa-rs/sdk/tests/fixtures/C.jpg \
    -o /tmp/signed.jpg \
    -c ../c2pa-rs/sdk/tests/fixtures/certs/ps256.pub \
    -k ../c2pa-rs/sdk/tests/fixtures/certs/ps256.pem
```

Both subcommands accept `-s settings.toml` to load a settings file; see
[c2pa-cli/settings.toml](c2pa-cli/settings.toml) for a working sample with a
custom trust list and verification policy.

## License

This project is dual-licensed under the MIT license and the Apache License 2.0
(matching c2pa-rs). See [LICENSE](LICENSE) for details.
