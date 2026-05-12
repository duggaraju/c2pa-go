# c2pa-go

Go bindings for the [c2pa-rs](https://github.com/contentauth/c2pa-rs) Content
Authenticity SDK. This package exposes the C2PA reader, builder, signer and
settings APIs through cgo so that Go programs can read, create and sign
[C2PA](https://c2pa.org) manifests.

The bindings track c2pa-rs `0.83.x` and link against the `c2pa_c` shared library
produced by the `c2pa_c_ffi` crate.

## Status

- Supported assets: anything supported by c2pa-rs (JPEG, PNG, TIFF, MP4, etc.).
- Supported algorithms: `Ps256`/`Ps384`/`Ps512`, `Es256`/`Es384`/`Es512`,
  `Ed25519`.
- Bindings are not yet API-stable; the surface may change before a tagged
  release.

## Repository layout

```
lib/        Go package (module github.com/duggaraju/c2pa-go/lib)
example/    Sample CLI that reads and signs assets
schema/     Generated Go types from the c2pa JSON schemas (Work in progress)
c2pa-rs/    Pinned upstream submodule built into the c2pa_c shared lib
build.sh    Convenience build script (cargo + go)
```

## Requirements

- Go 1.24+
- A C toolchain (`gcc`/`clang`) with `cgo` enabled
- Rust toolchain (only required to rebuild `c2pa-rs`)
- `git submodule update --init --recursive` for the pinned c2pa-rs checkout

## Building

```sh
git submodule update --init --recursive
./build.sh            # debug build
./build.sh release    # release build (-tags=release)
```

`build.sh` runs `cargo build --features rust_native_crypto,http` inside
`c2pa-rs/c2pa_c_ffi` and then `go build` for both `lib` and `example`.

In your own Go module, add the package as a dependency:

```sh
go get github.com/duggaraju/c2pa-go/lib
```

You will still need a built `libc2pa_c` available at link time. The cgo flag
files in [lib](lib) point at `c2pa-rs/target/{debug,release}`; adjust them for
your deployment if you ship the shared library elsewhere.

## Usage

All operations are scoped to a `*lib.Context`. A context optionally carries
`Settings` (trust list, verification policy, builder defaults, etc.).

### Reading a manifest

```go
package main

import (
    "fmt"
    "log"

    "github.com/duggaraju/c2pa-go/lib"
)

func main() {
    ctx, err := lib.NewContext()
    if err != nil {
        log.Fatal(err)
    }
    defer ctx.Close()

    r, err := lib.ReaderFromFile(ctx, "signed.jpg")
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
builder, err := lib.NewContextBuilder()
if err != nil { log.Fatal(err) }
defer builder.Close()

content, _ := os.ReadFile("settings.toml")
settings, err := lib.NewSettings()
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

Provide a manifest definition and an implementation of `lib.Signer`. The
binding will invoke `Sign` via a cgo callback when c2pa-rs needs a COSE
signature over the claim bytes.

```go
manifestJson := `{
  "claim_generator": "my-app/1.0",
  "title": "example.jpg"
}`

builder, err := lib.BuilderFromJson(ctx, manifestJson)
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

func (s *MyPs256Signer) Alg() lib.SigningAlg     { return lib.SigningAlgPs256 }
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
builder.SetIntent(lib.IntentCreate, lib.SourceDigitalCapture)
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
