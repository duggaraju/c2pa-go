package c2pa

/*
#include <c2pa.h>
*/
import "C"

import "unsafe"

// Version returns the version string from the c2pa library.
func Version() string {
	cs := C.c2pa_version()
	return C.GoString(cs)
}

func SetError(error string) {
	cstr := C.CString(error)
	defer C.free(unsafe.Pointer(cstr))
	C.c2pa_error_set_last(cstr)
}

func c2paError() string {
	cs := C.c2pa_error()
	defer C.c2pa_free(unsafe.Pointer(cs))
	return C.GoString(cs)
}

// BuilderSupportedMimeTypes returns the MIME types supported by the C2PA builder.
func BuilderSupportedMimeTypes() []string {
	var count C.uintptr_t
	arr := C.c2pa_builder_supported_mime_types(&count)
	if arr == nil || count == 0 {
		return nil
	}
	defer C.c2pa_free_string_array(arr, count)
	return cStringArrayToGo(arr, int(count))
}

// ReaderSupportedMimeTypes returns the MIME types supported by the C2PA reader.
func ReaderSupportedMimeTypes() []string {
	var count C.uintptr_t
	arr := C.c2pa_reader_supported_mime_types(&count)
	if arr == nil || count == 0 {
		return nil
	}
	defer C.c2pa_free_string_array(arr, count)
	return cStringArrayToGo(arr, int(count))
}

// cStringArrayToGo converts a C array of `count` C strings into a []string.
func cStringArrayToGo(arr **C.char, count int) []string {
	out := make([]string, count)
	slice := unsafe.Slice(arr, count)
	for i, p := range slice {
		out[i] = C.GoString(p)
	}
	return out
}
