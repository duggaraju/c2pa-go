//go:build release && unix
// +build release,unix

// Release build of the c2pa-rs C FFI static library (Unix-like systems).
//
// Run `go generate -tags=release ./c2pa/...` to build the release library
// into c2pa-rs/target/release, matching the paths in cgo_flags_release*.go.

package c2pa

//go:generate sh -c "cd ../c2pa-rs/c2pa_c_ffi && cargo build --release --no-default-features --features rust_native_crypto,http,add_thumbnails"
