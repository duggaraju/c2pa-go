//go:build !release && windows
// +build !release,windows

// Debug build of the c2pa-rs C FFI static library (Windows).
//
// Run `go generate ./c2pa/...` (no tags) to build the debug library into
// c2pa-rs/target/debug, matching the paths in cgo_flags_debug*.go.

package c2pa

//go:generate cmd /C "cd ../c2pa-rs/c2pa_c_ffi && cargo build --no-default-features --features rust_native_crypto,http,add_thumbnails"
