//go:build !cgo

// Pure-Go stubs used when the package is built without cgo. These let the
// package compile and have its documentation rendered by tools like
// pkg.go.dev which build with CGO_ENABLED=0. At runtime every constructor
// returns an error; nothing actually talks to the libc2pa_c library.

package c2pa

import "unsafe"

// nativeAvailable reports whether the package was built with cgo support.
const nativeAvailable = false

const nocgoErr = "c2pa: native library unavailable: package built without cgo"

// ---- helpers ---------------------------------------------------------------

func c2paFree(unsafe.Pointer)                       {}
func takeCBytes(unsafe.Pointer, int64) []byte    { return nil }
func c2paFreeStringArray(unsafe.Pointer, int)       {}

// ---- global ---------------------------------------------------------------

func c2paVersion() string             { return "" }
func c2paError() string           { return nocgoErr }
func c2paErrorSetLast(string)         {}
func c2paBuilderSupportedMimeTypes() []string { return nil }
func c2paReaderSupportedMimeTypes() []string  { return nil }

// ---- context ---------------------------------------------------------------

func c2paContextNew() unsafe.Pointer       { return nil }
func c2paContextCancel(unsafe.Pointer) int { return -1 }

// ---- context builder -------------------------------------------------------

func c2paContextBuilderNew() unsafe.Pointer                                         { return nil }
func c2paContextBuilderSetSigner(unsafe.Pointer, unsafe.Pointer) int                { return -1 }
func c2paContextBuilderSetSettings(unsafe.Pointer, unsafe.Pointer) int              { return -1 }
func c2paContextBuilderSetHttpResolver(unsafe.Pointer, unsafe.Pointer) int          { return -1 }
func c2paContextBuilderSetProgressCallback(unsafe.Pointer, uintptr) int             { return -1 }
func c2paContextBuilderBuild(unsafe.Pointer) unsafe.Pointer                         { return nil }

// ---- reader ----------------------------------------------------------------

func c2paReaderNew() unsafe.Pointer                                                                       { return nil }
func c2paReaderFromContext(unsafe.Pointer) unsafe.Pointer                                                  { return nil }
func c2paReaderJson(unsafe.Pointer) string                                                                 { return "" }
func c2paReaderDetailedJson(unsafe.Pointer) string                                                         { return "" }
func c2paReaderRemoteUrl(unsafe.Pointer) string                                                            { return "" }
func c2paReaderIsEmbedded(unsafe.Pointer) bool                                                             { return false }
func c2paReaderResourceToStream(unsafe.Pointer, string, unsafe.Pointer) int64                              { return -1 }
func c2paReaderWithStream(unsafe.Pointer, string, unsafe.Pointer) unsafe.Pointer                           { return nil }
func c2paReaderWithManifestDataAndStream(unsafe.Pointer, string, unsafe.Pointer, []byte) unsafe.Pointer             { return nil }
func c2paReaderWithFragment(unsafe.Pointer, unsafe.Pointer, unsafe.Pointer, string) unsafe.Pointer         { return nil }

// ---- builder ---------------------------------------------------------------

func c2paBuilderFromContext(unsafe.Pointer) unsafe.Pointer                                                   { return nil }
func c2paBuilderWithDefinition(unsafe.Pointer, string) unsafe.Pointer                                        { return nil }
func c2paBuilderWithArchive(unsafe.Pointer, unsafe.Pointer) unsafe.Pointer                                   { return nil }
func c2paBuilderSetNoEmbed(unsafe.Pointer)                                                                   {}
func c2paBuilderSetRemoteUrl(unsafe.Pointer, string) int                                                     { return -1 }
func c2paBuilderSetBasePath(unsafe.Pointer, string) int                                                      { return -1 }
func c2paBuilderSetIntent(unsafe.Pointer, BuilderIntent, DigitalSourceType) int                              { return -1 }
func c2paBuilderAddAction(unsafe.Pointer, string) int                                                        { return -1 }
func c2paBuilderAddResource(unsafe.Pointer, string, unsafe.Pointer) int                                      { return -1 }
func c2paBuilderAddIngredientFromStream(unsafe.Pointer, string, string, unsafe.Pointer) int                  { return -1 }
func c2paBuilderToArchive(unsafe.Pointer, unsafe.Pointer) int                                                { return -1 }
func c2paBuilderAddIngredientFromArchive(unsafe.Pointer, unsafe.Pointer) int                                 { return -1 }
func c2paBuilderWriteIngredientArchive(unsafe.Pointer, string, unsafe.Pointer) int                           { return -1 }
func c2paBuilderSign(unsafe.Pointer, string, unsafe.Pointer, unsafe.Pointer, unsafe.Pointer) ([]byte, int64) { return nil, -1 }
func c2paBuilderSignContext(unsafe.Pointer, string, unsafe.Pointer, unsafe.Pointer) ([]byte, int64)          { return nil, -1 }
func c2paBuilderNeedsPlaceholder(unsafe.Pointer, string) int                                                 { return -1 }
func c2paBuilderHashType(unsafe.Pointer, string) (HashType, int)                                             { return "", -1 }
func c2paBuilderPlaceholder(unsafe.Pointer, string) ([]byte, int64)                                          { return nil, -1 }
func c2paBuilderDataHashedPlaceholder(unsafe.Pointer, uintptr, string) ([]byte, int64)                       { return nil, -1 }
func c2paBuilderSignDataHashedEmbeddable(unsafe.Pointer, unsafe.Pointer, string, string, unsafe.Pointer) ([]byte, int64) {
	return nil, -1
}
func c2paBuilderSignEmbeddable(unsafe.Pointer, string) ([]byte, int64) { return nil, -1 }
func c2paBuilderSetDataHashExclusions(unsafe.Pointer, []uint64) int    { return -1 }
func c2paBuilderSetFixedSizeMerkle(unsafe.Pointer, uintptr) int        { return -1 }
func c2paBuilderHashMdatBytes(unsafe.Pointer, uintptr, []byte, bool) int {
	return -1
}
func c2paBuilderUpdateHashFromStream(unsafe.Pointer, string, unsafe.Pointer) int { return -1 }
func c2paFormatEmbeddable(string, []byte) ([]byte, int64)                        { return nil, -1 }

// ---- signer ----------------------------------------------------------------

func c2paSignerReserveSize(unsafe.Pointer) int64                          { return -1 }
func c2paSignerCreate(uintptr, SigningAlg, string, string) unsafe.Pointer { return nil }
func c2paSignerFromInfo(SigningAlg, string, string, string) unsafe.Pointer    { return nil }
func c2paIdentitySignerCreate(unsafe.Pointer, unsafe.Pointer, []string, []string) unsafe.Pointer {
	return nil
}
func c2paEd25519Sign([]byte, string) []byte { return nil }

// ---- http resolver ---------------------------------------------------------

func c2paHttpResolverCreate(uintptr) unsafe.Pointer { return nil }

// ---- settings --------------------------------------------------------------

func c2paSettingsNew() unsafe.Pointer                                  { return nil }
func c2paSettingsUpdateFromString(unsafe.Pointer, string, string) int  { return -1 }
func c2paSettingsSetValue(unsafe.Pointer, string, string) int          { return -1 }

// ---- stream ----------------------------------------------------------------

func c2paCreateStream(uintptr) unsafe.Pointer { return nil }
func c2paReleaseStream(unsafe.Pointer)        {}
