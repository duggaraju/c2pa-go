//go:build linux && !release
// +build linux,!release

package lib

/*
#cgo CFLAGS: -I${SRCDIR}/../c2pa-rs/target/debug
#cgo LDFLAGS: -L${SRCDIR}/../c2pa-rs/target/debug -Wl,-Bstatic -lc2pa_c -Wl,-Bdynamic -lm
*/
import "C"
