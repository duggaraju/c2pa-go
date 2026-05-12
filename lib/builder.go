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

// BuilderIntent corresponds to C2paBuilderIntent.
type BuilderIntent C.C2paBuilderIntent

const (
	IntentCreate BuilderIntent = C.Create
	IntentEdit   BuilderIntent = C.Edit
	IntentUpdate BuilderIntent = C.Update
)

// DigitalSourceType corresponds to C2paDigitalSourceType.
type DigitalSourceType C.C2paDigitalSourceType

const (
	SourceEmpty                                DigitalSourceType = C.Empty
	SourceTrainedAlgorithmicData               DigitalSourceType = C.TrainedAlgorithmicData
	SourceDigitalCapture                       DigitalSourceType = C.DigitalCapture
	SourceComputationalCapture                 DigitalSourceType = C.ComputationalCapture
	SourceNegativeFilm                         DigitalSourceType = C.NegativeFilm
	SourcePositiveFilm                         DigitalSourceType = C.PositiveFilm
	SourcePrint                                DigitalSourceType = C.Print
	SourceHumanEdits                           DigitalSourceType = C.HumanEdits
	SourceCompositeWithTrainedAlgorithmicMedia DigitalSourceType = C.CompositeWithTrainedAlgorithmicMedia
	SourceAlgorithmicallyEnhanced              DigitalSourceType = C.AlgorithmicallyEnhanced
	SourceDigitalCreation                      DigitalSourceType = C.DigitalCreation
	SourceDataDrivenMedia                      DigitalSourceType = C.DataDrivenMedia
	SourceTrainedAlgorithmicMedia              DigitalSourceType = C.TrainedAlgorithmicMedia
	SourceAlgorithmicMedia                     DigitalSourceType = C.AlgorithmicMedia
	SourceScreenCapture                        DigitalSourceType = C.ScreenCapture
	SourceVirtualRecording                     DigitalSourceType = C.VirtualRecording
	SourceComposite                            DigitalSourceType = C.Composite
	SourceCompositeCapture                     DigitalSourceType = C.CompositeCapture
	SourceCompositeSynthetic                   DigitalSourceType = C.CompositeSynthetic
)

func (b *Builder) Close() {
	C.c2pa_free(unsafe.Pointer(b.ptr))
}

func (b *Builder) SetNoEmbed() {
	C.c2pa_builder_set_no_embed(b.ptr)
}

// SetRemoteUrl sets the remote URL that will be embedded into the asset on signing.
func (b *Builder) SetRemoteUrl(url string) error {
	curl := C.CString(url)
	defer C.free(unsafe.Pointer(curl))
	if C.c2pa_builder_set_remote_url(b.ptr, curl) < 0 {
		return fmt.Errorf("failed to set remote url: %s", C2paError())
	}
	return nil
}

// SetBasePath sets the directory used to resolve resources not found in memory.
func (b *Builder) SetBasePath(path string) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	if C.c2pa_builder_set_base_path(b.ptr, cpath) < 0 {
		return fmt.Errorf("failed to set base path: %s", C2paError())
	}
	return nil
}

// SetIntent sets the builder intent. digitalSourceType is required for IntentCreate
// and ignored for other intents.
func (b *Builder) SetIntent(intent BuilderIntent, digitalSourceType DigitalSourceType) error {
	if C.c2pa_builder_set_intent(b.ptr, uint32(intent), uint32(digitalSourceType)) < 0 {
		return fmt.Errorf("failed to set intent: %s", C2paError())
	}
	return nil
}

// AddAction adds an action assertion described by the given JSON string.
func (b *Builder) AddAction(actionJson string) error {
	cjson := C.CString(actionJson)
	defer C.free(unsafe.Pointer(cjson))
	if C.c2pa_builder_add_action(b.ptr, cjson) < 0 {
		return fmt.Errorf("failed to add action: %s", C2paError())
	}
	return nil
}

// AddResource adds a resource read from file under the given URI identifier.
func (b *Builder) AddResource(uri string, file *os.File) error {
	curi := C.CString(uri)
	defer C.free(unsafe.Pointer(curi))

	stream, err := NewStream(file)
	if err != nil {
		return err
	}
	defer stream.Close()

	if C.c2pa_builder_add_resource(b.ptr, curi, stream.ptr) < 0 {
		return fmt.Errorf("failed to add resource %s: %s", uri, C2paError())
	}
	return nil
}

// AddResourceFromFile is a convenience wrapper that opens path and adds it as a resource.
func (b *Builder) AddResourceFromFile(uri string, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", path, err)
	}
	defer f.Close()
	return b.AddResource(uri, f)
}

// AddIngredientFromStream adds an ingredient described by ingredientJson, read
// from file with the given format (mime type or file extension).
func (b *Builder) AddIngredientFromStream(ingredientJson string, format string, file *os.File) error {
	cjson := C.CString(ingredientJson)
	defer C.free(unsafe.Pointer(cjson))
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))

	stream, err := NewStream(file)
	if err != nil {
		return err
	}
	defer stream.Close()

	if C.c2pa_builder_add_ingredient_from_stream(b.ptr, cjson, cformat, stream.ptr) < 0 {
		return fmt.Errorf("failed to add ingredient: %s", C2paError())
	}
	return nil
}

// AddIngredientFromFile is a convenience wrapper that opens path and adds it as
// an ingredient, deriving the format from the file extension.
func (b *Builder) AddIngredientFromFile(ingredientJson string, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", path, err)
	}
	defer f.Close()
	ext := filepath.Ext(path)
	if len(ext) > 0 {
		ext = ext[1:]
	}
	return b.AddIngredientFromStream(ingredientJson, ext, f)
}

// ToArchive writes a builder archive to the given file.
func (b *Builder) ToArchive(file *os.File) error {
	stream, err := NewStream(file)
	if err != nil {
		return err
	}
	defer stream.Close()

	if C.c2pa_builder_to_archive(b.ptr, stream.ptr) < 0 {
		return fmt.Errorf("failed to write archive: %s", C2paError())
	}
	return nil
}

// ToArchiveFile is a convenience wrapper that creates path and writes the
// builder archive to it.
func (b *Builder) ToArchiveFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create %s: %v", path, err)
	}
	defer f.Close()
	return b.ToArchive(f)
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

	var manifest *C.uchar
	len := C.sign_data(b.ptr, cformat, input_stream.ptr, output_stream.ptr, signerAdapter.ptr, unsafe.Pointer(&manifest))
	if len < 0 {
		return nil, fmt.Errorf("failed to sign file %s: %s", input_file, C2paError())
	}

	defer C.c2pa_free(unsafe.Pointer(manifest))
	return C.GoBytes(unsafe.Pointer(manifest), C.int(len)), nil
}

// NewBuilder creates a new Builder from the given Context.
func NewBuilder(ctx *Context) (*Builder, error) {
	if ctx == nil || ctx.ptr == nil {
		return nil, fmt.Errorf("context is nil")
	}
	ptr := C.c2pa_builder_from_context(ctx.ptr)
	if ptr == nil {
		return nil, fmt.Errorf("failed to create c2pa Builder: %s", C2paError())
	}
	return &Builder{ptr: ptr}, nil
}

// BuilderFromJson creates a Builder from the given JSON manifest definition
// using the supplied Context.
func BuilderFromJson(ctx *Context, json string) (*Builder, error) {
	b, err := NewBuilder(ctx)
	if err != nil {
		return nil, err
	}

	cjson := C.CString(json)
	defer C.free(unsafe.Pointer(cjson))

	ptr := C.c2pa_builder_with_definition(b.ptr, cjson)
	if ptr == nil {
		return nil, fmt.Errorf("failed to set definition on c2pa Builder: %s", C2paError())
	}
	b.ptr = ptr
	return b, nil
}

// BuilderFromArchive creates a Builder from an archive previously produced by
// ToArchive, using the supplied Context.
func BuilderFromArchive(ctx *Context, file *os.File) (*Builder, error) {
	b, err := NewBuilder(ctx)
	if err != nil {
		return nil, err
	}

	stream, err := NewStream(file)
	if err != nil {
		b.Close()
		return nil, err
	}
	defer stream.Close()

	ptr := C.c2pa_builder_with_archive(b.ptr, stream.ptr)
	if ptr == nil {
		b.Close()
		return nil, fmt.Errorf("failed to load archive: %s", C2paError())
	}
	b.ptr = ptr
	return b, nil
}

// BuilderFromArchiveFile is a convenience wrapper around BuilderFromArchive.
func BuilderFromArchiveFile(ctx *Context, path string) (*Builder, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", path, err)
	}
	defer f.Close()
	return BuilderFromArchive(ctx, f)
}
