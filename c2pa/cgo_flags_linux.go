//go:build cgo && linux
// +build cgo,linux

package c2pa

/*
#cgo !release CFLAGS: -I${SRCDIR}/../c2pa-rs/target/debug
#cgo release CFLAGS: -I${SRCDIR}/../c2pa-rs/target/release
#cgo !release LDFLAGS: -L${SRCDIR}/../c2pa-rs/target/debug
#cgo release LDFLAGS: -L${SRCDIR}/../c2pa-rs/target/release
#cgo LDFLAGS: -Wl,-Bstatic -lc2pa_c -Wl,-Bdynamic -lm
*/
import "C"
