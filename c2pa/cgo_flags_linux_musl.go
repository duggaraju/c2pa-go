//go:build cgo && linux && musl
// +build cgo,linux,musl

// Linux flags for musl libc, for building on a glibc host. Enabled with
// `-tags=musl` together with CC=musl-gcc. (On Alpine or another musl host the
// stock flags in cgo_flags_linux.go already work.)
//
// `CC=musl-gcc go generate -tags=musl[,release] ./c2pa` builds the native
// library for <arch>-unknown-linux-musl, so the artifacts live under
// c2pa-rs/target/<arch>-unknown-linux-musl/{debug,release} rather than
// c2pa-rs/target/{debug,release}. Only amd64 and arm64 are wired up here; for
// other architectures set CGO_CFLAGS/CGO_LDFLAGS (see the README).
//
// libunwind.a is staged next to libc2pa_c.a by the go:generate directives in
// generate_musl_*.go. Both archives are named explicitly (rather than
// -Wl,-Bstatic ... -Wl,-Bdynamic) so that `-extldflags "-static"` still yields
// a fully static executable.

package c2pa

/*
#cgo amd64,!release CFLAGS: -I${SRCDIR}/../c2pa-rs/target/x86_64-unknown-linux-musl/debug
#cgo amd64,release CFLAGS: -I${SRCDIR}/../c2pa-rs/target/x86_64-unknown-linux-musl/release
#cgo amd64,!release LDFLAGS: -L${SRCDIR}/../c2pa-rs/target/x86_64-unknown-linux-musl/debug
#cgo amd64,release LDFLAGS: -L${SRCDIR}/../c2pa-rs/target/x86_64-unknown-linux-musl/release
#cgo arm64,!release CFLAGS: -I${SRCDIR}/../c2pa-rs/target/aarch64-unknown-linux-musl/debug
#cgo arm64,release CFLAGS: -I${SRCDIR}/../c2pa-rs/target/aarch64-unknown-linux-musl/release
#cgo arm64,!release LDFLAGS: -L${SRCDIR}/../c2pa-rs/target/aarch64-unknown-linux-musl/debug
#cgo arm64,release LDFLAGS: -L${SRCDIR}/../c2pa-rs/target/aarch64-unknown-linux-musl/release
#cgo LDFLAGS: -l:libc2pa_c.a -l:libunwind.a -lm
*/
import "C"
