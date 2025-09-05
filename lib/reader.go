package lib

//#include <c2pa.h>
import "C"

import (
	"fmt"
	"os"
	"path/filepath"
	"unsafe"
)

// Reader wraps a C C2paReader*.
// It holds the underlying C pointer and provides a place to attach methods
// that operate on the C reader.
type Reader struct {
	ptr *C.C2paReader
}

// Ptr returns the underlying C pointer. Use carefully; callers must not keep
// the pointer past the lifetime of the Reader or the underlying C resource.
func (r *Reader) Ptr() *C.C2paReader {
	return r.ptr
}

func (r *Reader) Close() {
	C.c2pa_reader_free(r.ptr)
}

func (r *Reader) Json() string {
	json := C.c2pa_reader_json(r.ptr)
	defer C.c2pa_release_string(json)
	return C.GoString(json)
}

// ReaderFromFile creates a Reader by opening the given file path using the
// underlying C library. Returns an error if the reader could not be created.
func ReaderFromFile(path string) (*Reader, error) {
	ext := filepath.Ext(path)
	cformat := C.CString(ext[1:]) // skip the dot
	defer C.free(unsafe.Pointer(cformat))

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", path, err)
	}
	defer file.Close()

	stream, err := NewStream(file)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	reader := C.c2pa_reader_from_stream(cformat, stream.ptr)
	if reader == nil {
		return nil, fmt.Errorf("failed to create c2pa reader for %s: %s", path, C2paError())
	}
	return &Reader{ptr: reader}, nil
}
