package c2pa

// Run `go generate ./lib/...` to:
//
//  1. Regenerate the schema package (from c2pa-rs JSON schemas).
//  2. Build the c2pa-rs C FFI static library (debug + release) that the cgo
//     bindings link against. Output is written to c2pa-rs/target/{debug,release},
//     matching the paths referenced by the cgo_flags_*.go files.
//
//go:generate go generate github.com/duggaraju/c2pa-go/c2pa/schema
//go:generate sh -c "cd ../c2pa-rs/c2pa_c_ffi && cargo build --features rust_native_crypto,http"
//go:generate sh -c "cd ../c2pa-rs/c2pa_c_ffi && cargo build --release --features rust_native_crypto,http"
