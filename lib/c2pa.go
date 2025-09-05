package lib

/*
#include <c2pa.h>
*/
import "C"

// CpaVersion returns the version string from the c2pa library.
func CpaVersion() string {
	cs := C.c2pa_version()
	return C.GoString(cs)
}

func C2paError() string {
	cs := C.c2pa_error()
	defer C.c2pa_release_string(cs)
	return C.GoString(cs)
}
