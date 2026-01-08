//go:build !release && !linux
// +build !release,!linux

package lib

/*
#cgo CFLAGS: -I${SRCDIR}/../c2pa-rs/target/debug
#cgo LDFLAGS: -L${SRCDIR}/../c2pa-rs/target/debug -lc2pa_c
*/
import "C"
