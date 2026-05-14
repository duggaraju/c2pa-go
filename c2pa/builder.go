package c2pa

// #include "c2pa_helper.h"
import "C"

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"unsafe"

	schema "github.com/duggaraju/c2pa-go/c2pa/schema"
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
		return fmt.Errorf("failed to set remote url: %s", c2paError())
	}
	return nil
}

// SetBasePath sets the directory used to resolve resources not found in memory.
func (b *Builder) SetBasePath(path string) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	if C.c2pa_builder_set_base_path(b.ptr, cpath) < 0 {
		return fmt.Errorf("failed to set base path: %s", c2paError())
	}
	return nil
}

// SetIntent sets the builder intent. digitalSourceType is required for IntentCreate
// and ignored for other intents.
func (b *Builder) SetIntent(intent BuilderIntent, digitalSourceType DigitalSourceType) error {
	if C.c2pa_builder_set_intent(b.ptr, uint32(intent), uint32(digitalSourceType)) < 0 {
		return fmt.Errorf("failed to set intent: %s", c2paError())
	}
	return nil
}

// AddAction adds an action assertion described by the given JSON string.
func (b *Builder) AddAction(actionJson string) error {
	cjson := C.CString(actionJson)
	defer C.free(unsafe.Pointer(cjson))
	if C.c2pa_builder_add_action(b.ptr, cjson) < 0 {
		return fmt.Errorf("failed to add action: %s", c2paError())
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
		return fmt.Errorf("failed to add resource %s: %s", uri, c2paError())
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
		return fmt.Errorf("failed to add ingredient: %s", c2paError())
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
		return fmt.Errorf("failed to write archive: %s", c2paError())
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

// AddIngredientFromArchive loads an ingredient from a single-ingredient C2PA
// archive previously written by WriteIngredientArchive and adds it to this
// builder.
func (b *Builder) AddIngredientFromArchive(file *os.File) error {
	stream, err := NewStream(file)
	if err != nil {
		return err
	}
	defer stream.Close()

	if C.c2pa_builder_add_ingredient_from_archive(b.ptr, stream.ptr) < 0 {
		return fmt.Errorf("failed to add ingredient from archive: %s", c2paError())
	}
	return nil
}

// AddIngredientFromArchiveFile opens path and adds it to this builder via
// AddIngredientFromArchive.
func (b *Builder) AddIngredientFromArchiveFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", path, err)
	}
	defer f.Close()
	return b.AddIngredientFromArchive(f)
}

// WriteIngredientArchive writes a single-ingredient C2PA archive identified by
// ingredientId to file. Requires the generate_c2pa_archive builder setting to
// be enabled on the Context.
func (b *Builder) WriteIngredientArchive(ingredientId string, file *os.File) error {
	cid := C.CString(ingredientId)
	defer C.free(unsafe.Pointer(cid))

	stream, err := NewStream(file)
	if err != nil {
		return err
	}
	defer stream.Close()

	if C.c2pa_builder_write_ingredient_archive(b.ptr, cid, stream.ptr) < 0 {
		return fmt.Errorf("failed to write ingredient archive: %s", c2paError())
	}
	return nil
}

// WriteIngredientArchiveFile creates path and writes a single-ingredient C2PA
// archive to it.
func (b *Builder) WriteIngredientArchiveFile(ingredientId, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create %s: %v", path, err)
	}
	defer f.Close()
	return b.WriteIngredientArchive(ingredientId, f)
}

func (b *Builder) SignStream(format string, input *os.File, output *os.File, signer Signer) ([]byte, error) {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))

	input_stream, err := NewStream(input)
	if err != nil {
		return nil, err
	}
	defer input_stream.Close()

	output_stream, err := NewStream(output)
	if err != nil {
		return nil, err
	}
	defer output_stream.Close()

	var manifest *C.uchar
	var size C.int64_t
	if signer != nil {
		adapter, err := newSigner(signer)
		if err != nil {
			return nil, err
		}
		defer adapter.Close()
		size = C.c2pa_builder_sign(b.ptr, cformat, input_stream.ptr, output_stream.ptr, adapter.ptr, &manifest)

	} else {
		size = C.c2pa_builder_sign_context(b.ptr, cformat, input_stream.ptr, output_stream.ptr, &manifest)
	}
	if size < 0 {
		return nil, fmt.Errorf("failed to sign : %s", c2paError())
	}
	defer C.c2pa_free(unsafe.Pointer(manifest))
	return C.GoBytes(unsafe.Pointer(manifest), C.int(size)), nil
}

func (b *Builder) SignFile(input_file string, output_file string, signer Signer) ([]byte, error) {
	ext := filepath.Ext(input_file)
	format := ""
	if len(ext) > 0 {
		format = ext[1:] // skip the dot
	}

	input, err := os.Open(input_file)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", input_file, err)
	}
	defer input.Close()

	output, err := os.Create(output_file)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %v", output_file, err)
	}
	defer output.Close()

	return b.SignStream(format, input, output, signer)
}

func (b *Builder) Sign(input_file string, output_file string, signer Signer) ([]byte, error) {
	return b.SignFile(input_file, output_file, signer)
}

// Sign using the signer and http resolver from the context.
func (b *Builder) SignWithContext(input_file string, output_file string) ([]byte, error) {
	return b.SignFile(input_file, output_file, nil)
}

// NewBuilder creates a new Builder from the given Context.
func NewBuilder(ctx *Context) (*Builder, error) {
	if ctx == nil || ctx.ptr == nil {
		return nil, fmt.Errorf("context is nil")
	}
	ptr := C.c2pa_builder_from_context(ctx.ptr)
	if ptr == nil {
		return nil, fmt.Errorf("failed to create c2pa Builder: %s", c2paError())
	}
	return &Builder{ptr: ptr}, nil
}

func (b *Builder) WithDefinition(json string) (*Builder, error) {
	cjson := C.CString(json)
	defer C.free(unsafe.Pointer(cjson))
	b.ptr = C.c2pa_builder_with_definition(b.ptr, cjson)
	if b.ptr == nil {
		return nil, fmt.Errorf("Failed to set definition: %s", c2paError())
	}
	return b, nil
}

// BuilderFromArchive creates a Builder from an archive previously produced by
// ToArchive, using the supplied Context.
func (b *Builder) FromArchive(file *os.File) (*Builder, error) {
	stream, err := NewStream(file)
	if err != nil {
		b.Close()
		return nil, err
	}
	defer stream.Close()

	ptr := C.c2pa_builder_with_archive(b.ptr, stream.ptr)
	if ptr == nil {
		return nil, fmt.Errorf("failed to load archive: %s", c2paError())
	}
	b.ptr = ptr
	return b, nil
}

// BuilderFromArchiveFile is a convenience wrapper around BuilderFromArchive.
func (b *Builder) FromArchiveFile(path string) (*Builder, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", path, err)
	}
	defer f.Close()
	return b.FromArchive(f)
}

// BuilderFromDefinition creates a Builder from a typed ManifestDefinition.
// It marshals the definition to JSON and forwards it to BuilderFromJson.
func (b *Builder) WithManifestDefinition(def *schema.ManifestDefinition) (*Builder, error) {
	data, err := json.Marshal(def)
	if err != nil {
		return nil, err
	}
	return b.WithDefinition(string(data))
}

// AddActionTyped marshals action to JSON and adds it as an action assertion.
func (b *Builder) AddActionTyped(action any) error {
	data, err := json.Marshal(action)
	if err != nil {
		return err
	}
	return b.AddAction(string(data))
}

// HashType corresponds to C2paHashType — the hash binding type a Builder will
// produce for a given format in the embeddable signing workflow.
type HashType uint32

const (
	// HashTypeData uses placeholder + exclusions + hash + sign (JPEG, PNG, …).
	HashTypeData HashType = HashType(C.DataHash)
	// HashTypeBmff uses placeholder + hash + sign (MP4, AVIF, HEIF/HEIC).
	HashTypeBmff HashType = HashType(C.BmffHash)
	// HashTypeBox uses hash + sign with no placeholder needed.
	HashTypeBox HashType = HashType(C.BoxHash)
)

// HashExclusion describes a contiguous byte range to exclude from a DataHash
// binding.
type HashExclusion struct {
	Start  uint64
	Length uint64
}

// NeedsPlaceholder reports whether a placeholder manifest is required for the
// given format.
func (b *Builder) NeedsPlaceholder(format string) (bool, error) {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))
	rc := C.c2pa_builder_needs_placeholder(b.ptr, cformat)
	switch rc {
	case 1:
		return true, nil
	case 0:
		return false, nil
	default:
		return false, fmt.Errorf("failed to query placeholder need: %s", c2paError())
	}
}

// HashType returns the hash binding type the builder will use for the given
// format.
func (b *Builder) HashType(format string) (HashType, error) {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))
	var out uint32
	if C.c2pa_builder_hash_type(b.ptr, cformat, &out) < 0 {
		return 0, fmt.Errorf("failed to query hash type: %s", c2paError())
	}
	return HashType(out), nil
}

// Placeholder returns a composed placeholder manifest that can be embedded
// directly into an asset to reserve space for the final signed manifest. The
// signer is obtained from the Builder's Context.
func (b *Builder) Placeholder(format string) ([]byte, error) {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))
	var bytesPtr *C.uchar
	n := C.c2pa_builder_placeholder(b.ptr, cformat, &bytesPtr)
	if n < 0 {
		return nil, fmt.Errorf("failed to create placeholder: %s", c2paError())
	}
	return takeCBytes(bytesPtr, int64(n)), nil
}

// DataHashedPlaceholder reserves reservedSize bytes for a signature and
// returns the resulting placeholder manifest bytes.
func (b *Builder) DataHashedPlaceholder(reservedSize int, format string) ([]byte, error) {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))
	var bytesPtr *C.uchar
	n := C.c2pa_builder_data_hashed_placeholder(b.ptr, C.uintptr_t(reservedSize), cformat, &bytesPtr)
	if n < 0 {
		return nil, fmt.Errorf("failed to create data-hashed placeholder: %s", c2paError())
	}
	return takeCBytes(bytesPtr, int64(n)), nil
}

// SignDataHashedEmbeddable signs the manifest using the supplied signer and a
// pre-computed data hash JSON. asset may be nil if the hash JSON already
// contains the computed hash values.
func (b *Builder) SignDataHashedEmbeddable(signer Signer, dataHashJson string, format string, asset *os.File) ([]byte, error) {
	adapter, err := newSigner(signer)
	if err != nil {
		return nil, err
	}
	defer adapter.Close()

	cdh := C.CString(dataHashJson)
	defer C.free(unsafe.Pointer(cdh))
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))

	var assetPtr *C.C2paStream
	if asset != nil {
		s, err := NewStream(asset)
		if err != nil {
			return nil, err
		}
		defer s.Close()
		assetPtr = s.ptr
	}

	var bytesPtr *C.uchar
	n := C.c2pa_builder_sign_data_hashed_embeddable(b.ptr, adapter.ptr, cdh, cformat, assetPtr, &bytesPtr)
	if n < 0 {
		return nil, fmt.Errorf("failed to sign data hashed embeddable: %s", c2paError())
	}
	return takeCBytes(bytesPtr, int64(n)), nil
}

// SignEmbeddable signs the manifest and returns composed bytes ready for
// embedding. Operates in placeholder mode (after Placeholder) or direct mode
// (when the builder already has a valid hard-binding assertion). The signer is
// obtained from the Builder's Context.
func (b *Builder) SignEmbeddable(format string) ([]byte, error) {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))
	var bytesPtr *C.uchar
	n := C.c2pa_builder_sign_embeddable(b.ptr, cformat, &bytesPtr)
	if n < 0 {
		return nil, fmt.Errorf("failed to sign embeddable: %s", c2paError())
	}
	return takeCBytes(bytesPtr, int64(n)), nil
}

// SetDataHashExclusions registers the byte ranges where the composed
// placeholder was embedded in the asset. Must be called after Placeholder and
// before UpdateHashFromStream for DataHash workflows.
func (b *Builder) SetDataHashExclusions(exclusions []HashExclusion) error {
	flat := make([]C.uint64_t, 0, len(exclusions)*2)
	for _, e := range exclusions {
		flat = append(flat, C.uint64_t(e.Start), C.uint64_t(e.Length))
	}
	var ptr *C.uint64_t
	if len(flat) > 0 {
		ptr = &flat[0]
	}
	if C.c2pa_builder_set_data_hash_exclusions(b.ptr, ptr, C.uintptr_t(len(exclusions))) < 0 {
		return fmt.Errorf("failed to set data hash exclusions: %s", c2paError())
	}
	return nil
}

// SetFixedSizeMerkle configures the builder to hash fixed-size chunks of data,
// producing a Merkle tree per mdat. The unit is KB.
func (b *Builder) SetFixedSizeMerkle(fixedSizeKB uint) error {
	if C.c2pa_builder_set_fixed_size_merkle(b.ptr, C.uintptr_t(fixedSizeKB)) < 0 {
		return fmt.Errorf("failed to set fixed size merkle: %s", c2paError())
	}
	return nil
}

// HashMdatBytes accumulates leaf hashes for an mdat box. mdatId starts at 0
// and increments per mdat in the asset; data must be supplied in write order.
func (b *Builder) HashMdatBytes(mdatId uint, data []byte, largeSize bool) error {
	var dataPtr *C.uchar
	if len(data) > 0 {
		dataPtr = (*C.uchar)(unsafe.Pointer(&data[0]))
	}
	if C.c2pa_builder_hash_mdat_bytes(b.ptr, C.uintptr_t(mdatId), dataPtr, C.uintptr_t(len(data)), C.bool(largeSize)) < 0 {
		return fmt.Errorf("failed to hash mdat bytes: %s", c2paError())
	}
	return nil
}

// UpdateHashFromStream updates the builder's hard-binding assertion by hashing
// the asset stream. Detects DataHash, BmffHash, or BoxHash automatically.
func (b *Builder) UpdateHashFromStream(format string, file *os.File) error {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))

	stream, err := NewStream(file)
	if err != nil {
		return err
	}
	defer stream.Close()

	if C.c2pa_builder_update_hash_from_stream(b.ptr, cformat, stream.ptr) < 0 {
		return fmt.Errorf("failed to update hash from stream: %s", c2paError())
	}
	return nil
}

// FormatEmbeddable converts a raw application/c2pa manifest into an
// embeddable byte sequence for the given asset format.
func FormatEmbeddable(format string, manifestBytes []byte) ([]byte, error) {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))
	var inPtr *C.uchar
	if len(manifestBytes) > 0 {
		inPtr = (*C.uchar)(unsafe.Pointer(&manifestBytes[0]))
	}
	var outPtr *C.uchar
	n := C.c2pa_format_embeddable(cformat, inPtr, C.uintptr_t(len(manifestBytes)), &outPtr)
	if n < 0 {
		return nil, fmt.Errorf("failed to format embeddable: %s", c2paError())
	}
	return takeCBytes(outPtr, int64(n)), nil
}

// takeCBytes copies size bytes from a C-allocated buffer into a Go slice and
// frees the C buffer.
func takeCBytes(ptr *C.uchar, size int64) []byte {
	if ptr == nil || size <= 0 {
		return nil
	}
	out := C.GoBytes(unsafe.Pointer(ptr), C.int(size))
	C.c2pa_free(unsafe.Pointer(ptr))
	return out
}
