# c2pa-go (example)

This workspace contains a small Go library that wraps a C API (via cgo) which in turn can use a Rust-based implementation (`c2pa-rs`) as a submodule. It also contains a simple example CLI that demonstrates reading and (placeholder) signing operations.

Prerequisites
- Go 1.22+
- Rust toolchain (cargo)
- A C compiler (gcc/clang)

Quick steps

1. Initialize the workspace submodules (if you haven't already)

   git submodule update --init --recursive

   Or, if you have a remote for the Rust repo you want to use, you can add it to `lib/c2pa-rs` and initialize it.

2. Build the Rust library

   The repository includes a helper script: `lib/build_c2pa_rs.sh`.

   - Build debug (default):

     ./lib/build_c2pa_rs.sh

   - Build release:

     ./lib/build_c2pa_rs.sh release

   The script runs `cargo build` in `lib/c2pa-rs/c2pa_c_ffi`.

3. Build / run the Go example

   - Run in debug mode (uses the debug build paths):

     LD_LIBRARY_PATH=./c2pa-rs/target/debug go run ./example

   - Run using the release build (uses the release build paths):

     go run -tags release ./example

4. Advanced: build Rust automatically before invoking the Go tool with -toolexec

   You can use `go build` / `go run` `-toolexec` to run a wrapper that first builds the Rust code and then runs the actual tool. This is handy in CI or when you want a single command to ensure native artifacts are available.

   Example (build release Rust then run the example):

     go run -toolexec 'sh -c "./lib/build_c2pa_rs.sh release; exec \"$@\""' ./example

   Notes:
   - The wrapper string should end with `exec "$@"` so the original tool (compiler/linker) is invoked with its arguments after the script finishes.
   - Using `-toolexec` runs the wrapper for each underlying tool invocation; prefer this for CI or when you need reproducible build steps.

5. Example usage of the CLI

   - Read a file:

     go run ./example read -in path/to/file

   - Sign a file (placeholder — signing not implemented in this example):

     go run ./example sign -in path/to/file -out path/to/signed

   - Print version:

     go run ./example -v

Notes and next steps
- The Go package uses cgo to call a C API. The `lib` package includes fallbacks so you can develop without the Rust implementation present, but functionality will be limited.
- If you want me to wire up real signing support, add proper cleanup (Close methods), or change how the Rust submodule is fetched (clone vs submodule), say which behavior you prefer and I will update the repository.
