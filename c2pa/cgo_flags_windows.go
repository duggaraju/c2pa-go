//go:build cgo && windows
// +build cgo,windows

package c2pa

/*
#cgo !release CFLAGS: -I${SRCDIR}/../c2pa-rs/target/debug
#cgo release CFLAGS: -I${SRCDIR}/../c2pa-rs/target/release
#cgo !release LDFLAGS: -L${SRCDIR}/../c2pa-rs/target/debug
#cgo release LDFLAGS: -L${SRCDIR}/../c2pa-rs/target/release
#cgo LDFLAGS: -lc2pa_c -lws2_32 -luserenv -ladvapi32 -lncrypt -lcrypt32 -lbcrypt -lsecur32 -lntdll -lkernel32 -lole32 -loleaut32 -lpsapi -liphlpapi
*/
import "C"
