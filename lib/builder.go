package lib

// #include "c2pa_helper.h"
import "C"

import (
	"fmt"
	"os"
	"path/filepath"
	"unsafe"
)

// Builder wraps a C C2paBuilder*.
// It holds the underlying C pointer and provides a place to attach methods
// that operate on the C Builder.
type Builder struct {
	ptr *C.C2paBuilder
}

func (b *Builder) Close() {
	C.c2pa_builder_free(b.ptr)
}

func (b *Builder) SetNoEmbed() {
	C.c2pa_builder_set_no_embed(b.ptr)
}

func (b *Builder) Sign(input_file string, output_file string, signer Signer) ([]byte, error) {
	ext := filepath.Ext(input_file)
	cformat := C.CString(ext[1:]) // skip the dot
	defer C.free(unsafe.Pointer(cformat))

	input, err := os.Open(input_file)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", input_file, err)
	}
	defer input.Close()

	input_stream, err := NewStream(input)
	if err != nil {
		return nil, err
	}
	defer input_stream.Close()

	output, err := os.Create(output_file)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %v", output_file, err)
	}
	defer output.Close()

	output_stream, err := NewStream(output)
	if err != nil {
		return nil, err
	}
	defer output_stream.Close()

	signerAdapter, err := NewSigner(signer)
	if err != nil {
		return nil, err
	}
	defer signerAdapter.Close()

	var manifest C.uchar
	len := C.sign_data(b.ptr, cformat, input_stream.ptr, output_stream.ptr, signerAdapter.ptr, unsafe.Pointer(&manifest))
	if len < 0 {
		return nil, fmt.Errorf("failed to sign file %s: %s", input_file, C2paError())
	}

	defer C.c2pa_manifest_bytes_free(&manifest)
	return C.GoBytes(unsafe.Pointer(&manifest), C.int(len)), nil
}

// BuilderFromJson creates a Builder from the given JSON string using the
// underlying C library. Returns an error if the Builder could not be created.
func BuilderFromJson(json string) (*Builder, error) {
	cjson := C.CString(json)
	defer C.free(unsafe.Pointer(cjson))

	ptr := C.c2pa_builder_from_json(cjson)
	if ptr == nil {
		return nil, fmt.Errorf("failed to create c2pa Builder from JSON: %s", C2paError())
	}
	return &Builder{ptr: ptr}, nil
}
