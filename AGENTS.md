# AGENTS.md

This repository contains Go bindings for the upstream `c2pa-rs` project.
Work from the Go wrapper layer first, and only change the vendored Rust submodule when the task is explicitly about the FFI, schema export, or an upstream version bump.

## Primary Areas

- `c2pa/`: main Go module and cgo wrapper surface.
- `c2pa-cli/`: sample CLI in a separate Go module (`github.com/duggaraju/c2pa-go/c2pa-cli`), included via `go.work`.
- `c2pa/schema/`: generated Go types from upstream JSON schemas.
- `c2pa-rs/`: pinned upstream submodule; treat as external unless the task requires coordinated changes.

## Build And Test

Prefer the narrowest command that validates the slice you changed.

- Debug package tests: `cd /home/krishndu/go/c2pa-go && go test ./c2pa`
- Release package tests: `cd /home/krishndu/go/c2pa-go && go test -tags=release ./c2pa`
- Release builds: `cd /home/krishndu/go/c2pa-go && go build -tags=release ./c2pa && cd c2pa-cli && go build -tags=release .`
- Initialize the upstream checkout before build work: `git submodule update --init --recursive`

## Native Library Flow

The Go package links against prebuilt or locally built native artifacts under `c2pa-rs/target/`.

- Debug builds expect native artifacts in `c2pa-rs/target/debug`.
- Release builds expect native artifacts in `c2pa-rs/target/release`.
- To rebuild from source, use the repo hooks instead of ad hoc cargo commands when possible:
  - `go generate ./c2pa/...`
  - `go generate -tags=release ./c2pa/...`
- To use prebuilt release binaries, use `fetchlib`:
  - `go run ./c2pa/cmd/fetchlib`
  - `go run ./c2pa/cmd/fetchlib@latest`
  - Prebuilt Linux bundles are glibc only; `fetchlib` refuses to run on a musl host.

If the native library lives outside the default path, override `CGO_CFLAGS` and `CGO_LDFLAGS` rather than editing linker paths casually.

### musl

There are no prebuilt musl bundles; build locally.

- On Alpine/musl hosts the stock flow works unchanged (`go generate -tags=release ./c2pa`, then `go build -tags=release ./c2pa`).
- On a glibc host, cross-build with the `musl` tag: `export CC=musl-gcc`, then `go generate -tags=release,musl ./c2pa && go build -tags=release,musl ./c2pa`. Needs `musl-tools` and `rustup target add <arch>-unknown-linux-musl`.
- Test: `CC=musl-gcc go test -tags=release,musl ./c2pa`.
- Artifacts go to `c2pa-rs/target/<arch>-unknown-linux-musl/{debug,release}` (amd64 and arm64, host arch via `uname -m`). `generate_musl_*.go` build them with inline `go:generate` directives and copy Rust's `libunwind.a` alongside.
- Pass the same `-tags` to every `go` command. `cgo_flags_linux.go` is excluded by `musl`; `cgo_flags_linux_musl.go` replaces it.
- Keep the cargo feature list in `generate_musl_*.go` in sync with the other `generate_*` files.
- Full recipes (including other C libraries) are in README.md under "Building for musl and other C libraries".

## Schema Generation

`c2pa/schema/schema.go` is generated and checked in.
Only regenerate it when the task changes upstream schema output.

- Regenerate schema only: `cd /home/krishndu/go/c2pa-go && go generate ./c2pa/schema/...`

This path depends on Rust plus `node` and `npx quicktype`.

## Release Conventions

Keep these aligned when bumping versions:

- the `c2pa-rs` submodule revision
- `README.md`
- `c2pa/cmd/fetchlib/main.go` `DefaultTag`
- `c2pa-cli/go.mod` `require github.com/duggaraju/c2pa-go/c2pa vX.Y.Z`
- the Go module release tags

Use module-scoped tags in this repository. The `c2pa` and `c2pa-cli` modules
are released together at the same version, so a single bump produces two
tags off the same commit:

- `c2pa/vX.Y.Z`
- `c2pa-cli/vX.Y.Z`

Do not create a release tag from a dirty tree where the version bump is only in the working copy.
The CI workflow only attaches prebuilt `libc2pa_c` archives to the GitHub
Release produced by `c2pa/vX.Y.Z`; the `c2pa-cli/vX.Y.Z` tag only primes the
Go module proxy.

## Working Norms

- Prefer changing Go wrapper code in `c2pa/` over editing upstream Rust internals.
- Keep changes minimal and avoid reformatting generated files unless regeneration is part of the task.
- `c2pa/c2pa_test.go` depends on `linux`, `cgo`, and fixtures from the `c2pa-rs` submodule, so missing submodule content will break tests.
- When validating a change, start with `./c2pa` rather than `./...` unless the `c2pa-cli` module was also touched.

## Docs

Link to existing docs instead of duplicating them:

- [README.md](README.md): repo layout, build paths, prebuilt library flow, schema generation
- [c2pa-rs/docs/c_api.md](c2pa-rs/docs/c_api.md): upstream C API behind the Go bindings
- [c2pa-rs/docs/context-settings.md](c2pa-rs/docs/context-settings.md): settings model and behavior
- [c2pa-rs/docs/content_credentials.md](c2pa-rs/docs/content_credentials.md): domain concepts and terminology
- [c2pa-rs/docs/supported-formats.md](c2pa-rs/docs/supported-formats.md): supported asset formats
- [c2pa-rs/docs/release-process.md](c2pa-rs/docs/release-process.md): upstream release process details