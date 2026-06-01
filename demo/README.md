# Demo

This demo is intended to behave like an external consumer of the published `github.com/duggaraju/c2pa-go/c2pa` module. It builds a tiny Go program that links against the released `c2pa` package and prints the native C2PA library version returned by `c2pa.Version()`.

## Prerequisites

You need a working C toolchain with `cgo` enabled.

## Fetch the prebuilt native library

Download the release library into the location expected by the release cgo flags:

```sh
go run github.com/duggaraju/c2pa-go/c2pa/cmd/fetchlib
```

That command populates `c2pa-rs/target/release/` with `libc2pa_c` and headers for the current OS and architecture.

## Build the demo

Build the demo with the release tag so the Go wrapper links against the downloaded library:

```sh
go build -tags=release .
```

## Run the demo

Run the program from the repo root:

```sh
./demo
```

Expected output is a single line containing the C2PA library version, for example:

```text
c2pa version: <native-library-version>
```