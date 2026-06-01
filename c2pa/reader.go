package c2pa

//#include <c2pa.h>
import "C"

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"unsafe"

	schema "github.com/duggaraju/c2pa-go/c2pa/schema"
)

// Reader wraps a C C2paReader*.
// It holds the underlying C pointer and provides a place to attach methods
// that operate on the C reader.
type Reader struct {
	ptr *C.C2paReader
}

// NewDefaultReader creates a Reader using the library's default context.
func NewDefaultReader() (*Reader, error) {
	reader := C.c2pa_reader_new()
	if reader == nil {
		return nil, fmt.Errorf("failed to create default c2pa reader: %s", c2paError())
	}
	return &Reader{ptr: reader}, nil
}

func (r *Reader) Close() {
	C.c2pa_free(unsafe.Pointer(r.ptr))
}

func (r *Reader) Json() string {
	json := C.c2pa_reader_json(r.ptr)
	defer C.c2pa_free(unsafe.Pointer(json))
	return C.GoString(json)
}

// DetailedJson returns a detailed JSON description of the manifest store.
func (r *Reader) DetailedJson() string {
	json := C.c2pa_reader_detailed_json(r.ptr)
	defer C.c2pa_free(unsafe.Pointer(json))
	return C.GoString(json)
}

// RemoteUrl returns the remote URL the manifest was obtained from, or an empty
// string if the manifest was not remote.
func (r *Reader) RemoteUrl() string {
	s := C.c2pa_reader_remote_url(r.ptr)
	if s == nil {
		return ""
	}
	return C.GoString(s)
}

// IsEmbedded reports whether the reader was created from an embedded manifest.
func (r *Reader) IsEmbedded() bool {
	return bool(C.c2pa_reader_is_embedded(r.ptr))
}

// ResourceToStream writes the resource identified by uri to file and returns the
// number of bytes written.
func (r *Reader) ResourceToStream(uri string, file *os.File) (int64, error) {
	curi := C.CString(uri)
	defer C.free(unsafe.Pointer(curi))

	stream, err := NewStream(file)
	if err != nil {
		return 0, err
	}
	defer stream.Close()

	n := C.c2pa_reader_resource_to_stream(r.ptr, curi, stream.ptr)
	if n < 0 {
		return 0, fmt.Errorf("failed to write resource %s: %s", uri, c2paError())
	}
	return int64(n), nil
}

// ResourceToFile writes the resource identified by uri to the given path.
func (r *Reader) ResourceToFile(uri string, path string) (int64, error) {
	f, err := os.Create(path)
	if err != nil {
		return 0, fmt.Errorf("failed to create %s: %v", path, err)
	}
	defer f.Close()
	return r.ResourceToStream(uri, f)
}

// WithStream configures the reader from an asset stream. The reader must have
// been created with NewReader and not yet configured with another input.
func (r *Reader) WithStream(format string, file *os.File) error {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))

	stream, err := NewStream(file)
	if err != nil {
		return err
	}
	defer stream.Close()

	next := C.c2pa_reader_with_stream(r.ptr, cformat, stream.ptr)
	if next == nil {
		return fmt.Errorf("failed to configure reader with stream: %s", c2paError())
	}
	r.ptr = next
	return nil
}

// WithFile opens path, derives the format from its extension, and configures
// the reader via WithStream.
func (r *Reader) WithFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", path, err)
	}
	defer file.Close()

	format := filepath.Ext(path)
	if len(format) > 0 {
		format = format[1:]
	}

	return r.WithStream(format, file)
}

// WithManifestData configures the reader from an asset stream and a sidecar
// manifest. The reader must have been created with NewReader and not yet
// configured with WithStream/WithManifestData/WithFragment.
func (r *Reader) WithManifestData(format string, file *os.File, manifestData []byte) error {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))

	stream, err := NewStream(file)
	if err != nil {
		return err
	}
	defer stream.Close()

	var dataPtr *C.uchar
	if len(manifestData) > 0 {
		dataPtr = (*C.uchar)(unsafe.Pointer(&manifestData[0]))
	}
	next := C.c2pa_reader_with_manifest_data_and_stream(
		r.ptr, cformat, stream.ptr, dataPtr, C.uintptr_t(len(manifestData)))
	if next == nil {
		return fmt.Errorf("failed to configure reader with manifest data: %s", c2paError())
	}
	r.ptr = next
	return nil
}

// WithFragment configures the reader for a fragmented BMFF asset where the
// manifest lives in a separate fragment stream.
func (r *Reader) WithFragment(format string, asset *os.File, fragment *os.File) error {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))

	assetStream, err := NewStream(asset)
	if err != nil {
		return err
	}
	defer assetStream.Close()

	fragStream, err := NewStream(fragment)
	if err != nil {
		return err
	}
	defer fragStream.Close()

	next := C.c2pa_reader_with_fragment(r.ptr, cformat, assetStream.ptr, fragStream.ptr)
	if next == nil {
		return fmt.Errorf("failed to configure reader with fragment: %s", c2paError())
	}
	r.ptr = next
	return nil
}

// NewReader creates a new Reader from the given Context. The reader must be
// configured with a stream before it can be used.
func NewReader(ctx *Context) (*Reader, error) {
	if ctx == nil || ctx.ptr == nil {
		return nil, fmt.Errorf("context is nil")
	}
	reader := C.c2pa_reader_from_context(ctx.ptr)
	if reader == nil {
		return nil, fmt.Errorf("failed to create c2pa reader: %s", c2paError())
	}
	return &Reader{ptr: reader}, nil
}

// Manifest returns the manifest store parsed into the typed schema.ManifestStore
// representation.
func (r *Reader) Manifest() (*schema.ManifestStore, error) {
	var out schema.ManifestStore
	if err := json.Unmarshal([]byte(r.Json()), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DetailedManifest returns the detailed manifest store parsed into the typed
// schema.ManifestStore representation.
func (r *Reader) DetailedManifest() (*schema.ManifestStore, error) {
	var out schema.ManifestStore
	if err := json.Unmarshal([]byte(r.DetailedJson()), &out); err != nil {
		return nil, err
	}
	return &out, nil
}
