//go:build !release
// +build !release

// Debug build of the c2pa-rs C FFI static library.
//
// Run `go generate ./c2pa/...` (no tags) to build the debug library into
// c2pa-rs/target/debug, matching the paths in cgo_flags_debug*.go.

package c2pa

//go:generate sh -c "cd ../c2pa-rs/c2pa_c_ffi && cargo build --features rust_native_crypto"
