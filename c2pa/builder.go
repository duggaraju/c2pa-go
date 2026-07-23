package c2pa

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
	ptr unsafe.Pointer
}

func (b *Builder) Close() {
	c2paFree(b.ptr)
}

func (b *Builder) SetNoEmbed() {
	c2paBuilderSetNoEmbed(b.ptr)
}

// SetRemoteUrl sets the remote URL that will be embedded into the asset on signing.
func (b *Builder) SetRemoteUrl(url string) error {
	if c2paBuilderSetRemoteUrl(b.ptr, url) < 0 {
		return fmt.Errorf("failed to set remote url: %s", c2paError())
	}
	return nil
}

// SetBasePath sets the directory used to resolve resources not found in memory.
func (b *Builder) SetBasePath(path string) error {
	if c2paBuilderSetBasePath(b.ptr, path) < 0 {
		return fmt.Errorf("failed to set base path: %s", c2paError())
	}
	return nil
}

// SetIntent sets the builder intent. digitalSourceType is required for IntentCreate
// and ignored for other intents.
func (b *Builder) SetIntent(intent BuilderIntent, digitalSourceType DigitalSourceType) error {
	if c2paBuilderSetIntent(b.ptr, intent, digitalSourceType) < 0 {
		return fmt.Errorf("failed to set intent: %s", c2paError())
	}
	return nil
}

// AddAction adds an action assertion described by the given JSON string.
func (b *Builder) AddAction(actionJson string) error {
	if c2paBuilderAddAction(b.ptr, actionJson) < 0 {
		return fmt.Errorf("failed to add action: %s", c2paError())
	}
	return nil
}

// AddResource adds a resource read from file under the given URI identifier.
func (b *Builder) AddResource(uri string, file *os.File) error {
	stream, err := NewStream(file)
	if err != nil {
		return err
	}
	defer stream.Close()

	if c2paBuilderAddResource(b.ptr, uri, stream.ptr) < 0 {
		return fmt.Errorf("failed to add resource %s: %s", uri, c2paError())
	}
	return nil
}

// AddResourceFromFile is a convenience wrapper that opens path and adds it as a resource.
func (b *Builder) AddResourceFromFile(uri string, path string) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", path, err)
	}
	defer closeAndJoin(&err, f)
	return b.AddResource(uri, f)
}

// AddIngredientFromStream adds an ingredient described by ingredientJson, read
// from file with the given format (mime type or file extension).
func (b *Builder) AddIngredientFromStream(ingredientJson string, format string, file *os.File) error {
	stream, err := NewStream(file)
	if err != nil {
		return err
	}
	defer stream.Close()

	if c2paBuilderAddIngredientFromStream(b.ptr, ingredientJson, format, stream.ptr) < 0 {
		return fmt.Errorf("failed to add ingredient: %s", c2paError())
	}
	return nil
}

// AddIngredientFromFile is a convenience wrapper that opens path and adds it as
// an ingredient, deriving the format from the file extension.
func (b *Builder) AddIngredientFromFile(ingredientJson string, path string) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", path, err)
	}
	defer closeAndJoin(&err, f)
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

	if c2paBuilderToArchive(b.ptr, stream.ptr) < 0 {
		return fmt.Errorf("failed to write archive: %s", c2paError())
	}
	return nil
}

// ToArchiveFile is a convenience wrapper that creates path and writes the
// builder archive to it.
func (b *Builder) ToArchiveFile(path string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create %s: %v", path, err)
	}
	defer closeAndJoin(&err, f)
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

	if c2paBuilderAddIngredientFromArchive(b.ptr, stream.ptr) < 0 {
		return fmt.Errorf("failed to add ingredient from archive: %s", c2paError())
	}
	return nil
}

// AddIngredientFromArchiveFile opens path and adds it to this builder via
// AddIngredientFromArchive.
func (b *Builder) AddIngredientFromArchiveFile(path string) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", path, err)
	}
	defer closeAndJoin(&err, f)
	return b.AddIngredientFromArchive(f)
}

// WriteIngredientArchive writes a single-ingredient C2PA archive identified by
// ingredientId to file. Requires the generate_c2pa_archive builder setting to
// be enabled on the Context.
func (b *Builder) WriteIngredientArchive(ingredientId string, file *os.File) error {
	stream, err := NewStream(file)
	if err != nil {
		return err
	}
	defer stream.Close()

	if c2paBuilderWriteIngredientArchive(b.ptr, ingredientId, stream.ptr) < 0 {
		return fmt.Errorf("failed to write ingredient archive: %s", c2paError())
	}
	return nil
}

// WriteIngredientArchiveFile creates path and writes a single-ingredient C2PA
// archive to it.
func (b *Builder) WriteIngredientArchiveFile(ingredientId, path string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create %s: %v", path, err)
	}
	defer closeAndJoin(&err, f)
	return b.WriteIngredientArchive(ingredientId, f)
}

func (b *Builder) SignStream(format string, input *os.File, output *os.File, signer Signer) ([]byte, error) {
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

	var manifest []byte
	var n int64
	if signer != nil {
		native, err := takeNativeSigner(signer)
		if err != nil {
			return nil, err
		}
		defer native.Close()
		manifest, n = c2paBuilderSign(b.ptr, format, input_stream.ptr, output_stream.ptr, native.ptr)
	} else {
		manifest, n = c2paBuilderSignContext(b.ptr, format, input_stream.ptr, output_stream.ptr)
	}
	if n < 0 {
		return nil, fmt.Errorf("failed to sign : %s", c2paError())
	}
	return manifest, nil
}

func (b *Builder) SignFile(input_file string, output_file string, signer Signer) (manifest []byte, err error) {
	ext := filepath.Ext(input_file)
	format := ""
	if len(ext) > 0 {
		format = ext[1:] // skip the dot
	}

	input, err := os.Open(input_file)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", input_file, err)
	}
	defer closeAndJoin(&err, input)

	output, err := os.Create(output_file)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %v", output_file, err)
	}
	defer closeAndJoin(&err, output)

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
	ptr := c2paBuilderFromContext(ctx.ptr)
	if ptr == nil {
		return nil, fmt.Errorf("failed to create c2pa Builder: %s", c2paError())
	}
	return &Builder{ptr: ptr}, nil
}

func (b *Builder) WithDefinition(json string) (*Builder, error) {
	ptr := c2paBuilderWithDefinition(b.ptr, json)
	if ptr == nil {
		return nil, fmt.Errorf("failed to set definition: %s", c2paError())
	}
	b.ptr = ptr
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

	ptr := c2paBuilderWithArchive(b.ptr, stream.ptr)
	if ptr == nil {
		return nil, fmt.Errorf("failed to load archive: %s", c2paError())
	}
	b.ptr = ptr
	return b, nil
}

// BuilderFromArchiveFile is a convenience wrapper around BuilderFromArchive.
func (b *Builder) FromArchiveFile(path string) (_ *Builder, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", path, err)
	}
	defer closeAndJoin(&err, f)
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

// HashExclusion describes a contiguous byte range to exclude from a DataHash
// binding.
type HashExclusion struct {
	Start  uint64
	Length uint64
}

// NeedsPlaceholder reports whether a placeholder manifest is required for the
// given format.
func (b *Builder) NeedsPlaceholder(format string) (bool, error) {
	switch c2paBuilderNeedsPlaceholder(b.ptr, format) {
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
	out, rc := c2paBuilderHashType(b.ptr, format)
	if rc < 0 {
		return "", fmt.Errorf("failed to query hash type: %s", c2paError())
	}
	return out, nil
}

// Placeholder returns a composed placeholder manifest that can be embedded
// directly into an asset to reserve space for the final signed manifest. The
// signer is obtained from the Builder's Context.
func (b *Builder) Placeholder(format string) ([]byte, error) {
	out, n := c2paBuilderPlaceholder(b.ptr, format)
	if n < 0 {
		return nil, fmt.Errorf("failed to create placeholder: %s", c2paError())
	}
	return out, nil
}

// DataHashedPlaceholder reserves reservedSize bytes for a signature and
// returns the resulting placeholder manifest bytes.
func (b *Builder) DataHashedPlaceholder(reservedSize int, format string) ([]byte, error) {
	out, n := c2paBuilderDataHashedPlaceholder(b.ptr, uintptr(reservedSize), format)
	if n < 0 {
		return nil, fmt.Errorf("failed to create data-hashed placeholder: %s", c2paError())
	}
	return out, nil
}

// SignDataHashedEmbeddable signs the manifest using the supplied signer and a
// pre-computed data hash JSON. asset may be nil if the hash JSON already
// contains the computed hash values.
func (b *Builder) SignDataHashedEmbeddable(signer Signer, dataHashJson string, format string, asset *os.File) ([]byte, error) {
	native, err := takeNativeSigner(signer)
	if err != nil {
		return nil, err
	}
	defer native.Close()

	var assetPtr unsafe.Pointer
	if asset != nil {
		s, err := NewStream(asset)
		if err != nil {
			return nil, err
		}
		defer s.Close()
		assetPtr = s.ptr
	}

	out, n := c2paBuilderSignDataHashedEmbeddable(b.ptr, native.ptr, dataHashJson, format, assetPtr)
	if n < 0 {
		return nil, fmt.Errorf("failed to sign data hashed embeddable: %s", c2paError())
	}
	return out, nil
}

// SignEmbeddable signs the manifest and returns composed bytes ready for
// embedding. Operates in placeholder mode (after Placeholder) or direct mode
// (when the builder already has a valid hard-binding assertion). The signer is
// obtained from the Builder's Context.
func (b *Builder) SignEmbeddable(format string) ([]byte, error) {
	out, n := c2paBuilderSignEmbeddable(b.ptr, format)
	if n < 0 {
		return nil, fmt.Errorf("failed to sign embeddable: %s", c2paError())
	}
	return out, nil
}

// SetDataHashExclusions registers the byte ranges where the composed
// placeholder was embedded in the asset. Must be called after Placeholder and
// before UpdateHashFromStream for DataHash workflows.
func (b *Builder) SetDataHashExclusions(exclusions []HashExclusion) error {
	flat := make([]uint64, 0, len(exclusions)*2)
	for _, e := range exclusions {
		flat = append(flat, e.Start, e.Length)
	}
	if c2paBuilderSetDataHashExclusions(b.ptr, flat) < 0 {
		return fmt.Errorf("failed to set data hash exclusions: %s", c2paError())
	}
	return nil
}

// SetFixedSizeMerkle configures the builder to hash fixed-size chunks of data,
// producing a Merkle tree per mdat. The unit is KB.
func (b *Builder) SetFixedSizeMerkle(fixedSizeKB uint) error {
	if c2paBuilderSetFixedSizeMerkle(b.ptr, uintptr(fixedSizeKB)) < 0 {
		return fmt.Errorf("failed to set fixed size merkle: %s", c2paError())
	}
	return nil
}

// HashMdatBytes accumulates leaf hashes for an mdat box. mdatId starts at 0
// and increments per mdat in the asset; data must be supplied in write order.
func (b *Builder) HashMdatBytes(mdatId uint, data []byte, largeSize bool) error {
	if c2paBuilderHashMdatBytes(b.ptr, uintptr(mdatId), data, largeSize) < 0 {
		return fmt.Errorf("failed to hash mdat bytes: %s", c2paError())
	}
	return nil
}

// UpdateHashFromStream updates the builder's hard-binding assertion by hashing
// the asset stream. Detects DataHash, BmffHash, or BoxHash automatically.
func (b *Builder) UpdateHashFromStream(format string, file *os.File) error {
	stream, err := NewStream(file)
	if err != nil {
		return err
	}
	defer stream.Close()

	if c2paBuilderUpdateHashFromStream(b.ptr, format, stream.ptr) < 0 {
		return fmt.Errorf("failed to update hash from stream: %s", c2paError())
	}
	return nil
}

// FormatEmbeddable converts a raw application/c2pa manifest into an
// embeddable byte sequence for the given asset format.
func FormatEmbeddable(format string, manifestBytes []byte) ([]byte, error) {
	out, n := c2paFormatEmbeddable(format, manifestBytes)
	if n < 0 {
		return nil, fmt.Errorf("failed to format embeddable: %s", c2paError())
	}
	return out, nil
}
