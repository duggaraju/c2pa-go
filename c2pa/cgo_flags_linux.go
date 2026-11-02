//go:build cgo && linux && !musl
// +build cgo,linux,!musl

// Default Linux flags (glibc, or a musl host such as Alpine). They expect the
// native library built with a plain `cargo build` in c2pa-rs/target/{debug,release}.
// To cross-build for musl from a glibc host, use `-tags=musl`
// (cgo_flags_linux_musl.go).

package c2pa

/*
#cgo !release CFLAGS: -I${SRCDIR}/../c2pa-rs/target/debug
#cgo release CFLAGS: -I${SRCDIR}/../c2pa-rs/target/release
#cgo !release LDFLAGS: -L${SRCDIR}/../c2pa-rs/target/debug
#cgo release LDFLAGS: -L${SRCDIR}/../c2pa-rs/target/release
#cgo LDFLAGS: -Wl,-Bstatic -lc2pa_c -Wl,-Bdynamic -lm
*/
import "C"
