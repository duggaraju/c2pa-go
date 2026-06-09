//go:build cgo && linux && release
// +build cgo,linux,release

package c2pa

/*
#cgo CFLAGS: -I${SRCDIR}/../c2pa-rs/target/release
#cgo LDFLAGS: -L${SRCDIR}/../c2pa-rs/target/release -Wl,-Bstatic -lc2pa_c -Wl,-Bdynamic -lm
*/
import "C"
