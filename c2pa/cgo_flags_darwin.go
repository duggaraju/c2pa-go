//go:build cgo && darwin
// +build cgo,darwin

package c2pa

/*
#cgo !release CFLAGS: -I${SRCDIR}/../c2pa-rs/target/debug
#cgo release CFLAGS: -I${SRCDIR}/../c2pa-rs/target/release
#cgo !release LDFLAGS: -L${SRCDIR}/../c2pa-rs/target/debug
#cgo release LDFLAGS: -L${SRCDIR}/../c2pa-rs/target/release
#cgo LDFLAGS: -Wl,-search_paths_first -lc2pa_c -framework Security -framework CoreFoundation -framework SystemConfiguration -lresolv -ldl -lm
*/
import "C"
