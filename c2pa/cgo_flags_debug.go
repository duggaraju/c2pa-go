//go:build !release && !linux
// +build !release,!linux

package c2pa

/*
#cgo CFLAGS: -I${SRCDIR}/../c2pa-rs/target/debug
#cgo LDFLAGS: -L${SRCDIR}/../c2pa-rs/target/debug -lc2pa_c
*/
import "C"
