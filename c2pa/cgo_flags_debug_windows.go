//go:build windows && !release
// +build windows,!release

package c2pa

/*
#cgo CFLAGS: -I${SRCDIR}/../c2pa-rs/target/debug
#cgo LDFLAGS: -L${SRCDIR}/../c2pa-rs/target/debug -lc2pa_c -lws2_32 -luserenv -ladvapi32 -lncrypt -lcrypt32 -lbcrypt -lsecur32 -lntdll -lkernel32 -lole32 -loleaut32 -lpsapi -liphlpapi
*/
import "C"
