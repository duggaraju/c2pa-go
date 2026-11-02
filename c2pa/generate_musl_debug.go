//go:build !release && linux && musl
// +build !release,linux,musl

// Debug build of the c2pa-rs C FFI static library for musl libc, for a glibc
// host cross-building to musl (on Alpine/musl hosts the stock flow works).
//
// Run `CC=musl-gcc go generate -tags=musl ./c2pa/...` to build the debug
// library into c2pa-rs/target/<arch>-unknown-linux-musl/debug, matching the
// paths in cgo_flags_linux_musl.go. <arch> is the host's `uname -m`. The second
// directive stages the unwinder that the static library needs but does not
// bundle (glibc's libgcc_eh.a cannot be used with musl).
// Requires `rustup target add <arch>-unknown-linux-musl`.

package c2pa

//go:generate sh -c "cd ../c2pa-rs/c2pa_c_ffi && cargo build --target $(uname -m)-unknown-linux-musl --no-default-features --features rust_native_crypto,http,add_thumbnails"
//go:generate sh -c "cp $(rustc --print sysroot)/lib/rustlib/$(uname -m)-unknown-linux-musl/lib/self-contained/libunwind.a ../c2pa-rs/target/$(uname -m)-unknown-linux-musl/debug/"
