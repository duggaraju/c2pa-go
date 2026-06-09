//go:build cgo && darwin && !release
// +build cgo,darwin,!release

package c2pa

/*
#cgo CFLAGS: -I${SRCDIR}/../c2pa-rs/target/debug
#cgo LDFLAGS: -L${SRCDIR}/../c2pa-rs/target/debug -Wl,-search_paths_first -lc2pa_c -framework Security -framework CoreFoundation -framework SystemConfiguration -lresolv -ldl -lm
*/
import "C"
