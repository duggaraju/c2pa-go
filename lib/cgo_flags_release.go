//go:build release
// +build release

package lib

/*
#cgo CFLAGS: -I${SRCDIR}/../c2pa-rs/target/release
#cgo LDFLAGS: -L${SRCDIR}/../c2pa-rs/target/release -lc2pa_c
*/
import "C"
