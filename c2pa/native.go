//go:build cgo

// Package-internal cgo wrappers. All `import "C"` and direct calls into the
// libc2pa_c FFI live in this file so that the rest of the package can be
// compiled (and indexed by pkg.go.dev) without cgo. native_nocgo.go provides
// matching stubs for builds without cgo.

package c2pa

// #include <stdlib.h>
// #include <string.h>
// #include "c2pa_helper.h"
import "C"

import (
	"unsafe"
)

// nativeAvailable reports whether the package was built with cgo support.
const nativeAvailable = true

// ---- enum mappings ---------------------------------------------------------

// These maps translate the public string-typed enums (declared in
// constants.go) to the underlying libc2pa_c integer enums. Keeping the
// mapping here means upstream enum reordering can't silently break the Go
// public API.

var (
	signingAlgToC = map[SigningAlg]C.C2paSigningAlg{
		SigningAlgEs256:   C.Es256,
		SigningAlgEs384:   C.Es384,
		SigningAlgEs512:   C.Es512,
		SigningAlgPs256:   C.Ps256,
		SigningAlgPs384:   C.Ps384,
		SigningAlgPs512:   C.Ps512,
		SigningAlgEd25519: C.Ed25519,
	}

	intentToC = map[BuilderIntent]C.C2paBuilderIntent{
		IntentCreate: C.Create,
		IntentEdit:   C.Edit,
		IntentUpdate: C.Update,
	}

	sourceToC = map[DigitalSourceType]C.C2paDigitalSourceType{
		SourceEmpty:                                C.Empty,
		SourceTrainedAlgorithmicData:               C.TrainedAlgorithmicData,
		SourceDigitalCapture:                       C.DigitalCapture,
		SourceComputationalCapture:                 C.ComputationalCapture,
		SourceNegativeFilm:                         C.NegativeFilm,
		SourcePositiveFilm:                         C.PositiveFilm,
		SourcePrint:                                C.Print,
		SourceHumanEdits:                           C.HumanEdits,
		SourceCompositeWithTrainedAlgorithmicMedia: C.CompositeWithTrainedAlgorithmicMedia,
		SourceAlgorithmicallyEnhanced:              C.AlgorithmicallyEnhanced,
		SourceDigitalCreation:                      C.DigitalCreation,
		SourceDataDrivenMedia:                      C.DataDrivenMedia,
		SourceTrainedAlgorithmicMedia:              C.TrainedAlgorithmicMedia,
		SourceAlgorithmicMedia:                     C.AlgorithmicMedia,
		SourceScreenCapture:                        C.ScreenCapture,
		SourceVirtualRecording:                     C.VirtualRecording,
		SourceComposite:                            C.Composite,
		SourceCompositeCapture:                     C.CompositeCapture,
		SourceCompositeSynthetic:                   C.CompositeSynthetic,
	}

	hashTypeFromC = map[C.C2paHashType]HashType{
		C.DataHash: HashTypeData,
		C.BmffHash: HashTypeBmff,
		C.BoxHash:  HashTypeBox,
	}

	progressPhaseFromC = map[C.C2paProgressPhase]ProgressPhase{
		C.Reading:                ProgressReading,
		C.VerifyingManifest:      ProgressVerifyingManifest,
		C.VerifyingSignature:     ProgressVerifyingSignature,
		C.VerifyingIngredient:    ProgressVerifyingIngredient,
		C.VerifyingAssetHash:     ProgressVerifyingAssetHash,
		C.AddingIngredient:       ProgressAddingIngredient,
		C.Thumbnail:              ProgressThumbnail,
		C.Hashing:                ProgressHashing,
		C.Signing:                ProgressSigning,
		C.Embedding:              ProgressEmbedding,
		C.FetchingRemoteManifest: ProgressFetchingRemoteManifest,
		C.Writing:                ProgressWriting,
		C.FetchingOCSP:           ProgressFetchingOCSP,
		C.FetchingTimestamp:      ProgressFetchingTimestamp,
	}
)

// ---- helpers ---------------------------------------------------------------

func c2paFree(p unsafe.Pointer) {
	if p != nil {
		C.c2pa_free(p)
	}
}

// takeCBytes copies size bytes from a C-allocated buffer into a Go slice and
// frees the C buffer with c2pa_free.
func takeCBytes(p unsafe.Pointer, size int64) []byte {
	if p == nil || size <= 0 {
		return nil
	}
	out := C.GoBytes(p, C.int(size))
	C.c2pa_free(p)
	return out
}

func c2paFreeStringArray(arr unsafe.Pointer, count int) {
	if arr == nil || count <= 0 {
		return
	}
	C.c2pa_free_string_array((**C.char)(arr), C.uintptr_t(count))
}

// ---- global ---------------------------------------------------------------

func c2paVersion() string {
	cs := C.c2pa_version()
	defer C.c2pa_free(unsafe.Pointer(cs))
	return C.GoString(cs)
}

func c2paError() string {
	cs := C.c2pa_error()
	defer C.c2pa_free(unsafe.Pointer(cs))
	return C.GoString(cs)
}

func c2paErrorSetLast(msg string) {
	cs := C.CString(msg)
	defer C.free(unsafe.Pointer(cs))
	C.c2pa_error_set_last(cs)
}

func c2paBuilderSupportedMimeTypes() []string {
	var count C.uintptr_t
	arr := C.c2pa_builder_supported_mime_types(&count)
	if arr == nil || count == 0 {
		return nil
	}
	defer c2paFreeStringArray(unsafe.Pointer(arr), int(count))
	return cStringArrayToGo(arr, int(count))
}

func c2paReaderSupportedMimeTypes() []string {
	var count C.uintptr_t
	arr := C.c2pa_reader_supported_mime_types(&count)
	if arr == nil || count == 0 {
		return nil
	}
	defer c2paFreeStringArray(unsafe.Pointer(arr), int(count))
	return cStringArrayToGo(arr, int(count))
}

func cStringArrayToGo(arr **C.char, count int) []string {
	out := make([]string, count)
	slice := unsafe.Slice(arr, count)
	for i, p := range slice {
		out[i] = C.GoString(p)
	}
	return out
}

// ---- context ---------------------------------------------------------------

func c2paContextNew() unsafe.Pointer {
	return unsafe.Pointer(C.c2pa_context_new())
}

func c2paContextCancel(ctx unsafe.Pointer) int {
	return int(C.c2pa_context_cancel((*C.C2paContext)(ctx)))
}

// ---- context builder -------------------------------------------------------

func c2paContextBuilderNew() unsafe.Pointer {
	return unsafe.Pointer(C.c2pa_context_builder_new())
}

func c2paContextBuilderSetSigner(b, signer unsafe.Pointer) int {
	return int(C.c2pa_context_builder_set_signer((*C.C2paContextBuilder)(b), (*C.C2paSigner)(signer)))
}

func c2paContextBuilderSetSettings(b, settings unsafe.Pointer) int {
	return int(C.c2pa_context_builder_set_settings((*C.C2paContextBuilder)(b), (*C.C2paSettings)(settings)))
}

func c2paContextBuilderSetHttpResolver(b, resolver unsafe.Pointer) int {
	return int(C.c2pa_context_builder_set_http_resolver((*C.C2paContextBuilder)(b), (*C.C2paHttpResolver)(resolver)))
}

func c2paContextBuilderSetProgressCallback(b unsafe.Pointer, handle uintptr) int {
	return int(C.set_progress_callback((*C.C2paContextBuilder)(b), C.uintptr_t(handle)))
}

func c2paContextBuilderBuild(b unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.c2pa_context_builder_build((*C.C2paContextBuilder)(b)))
}

// ---- reader ----------------------------------------------------------------

func c2paReaderNew() unsafe.Pointer {
	return unsafe.Pointer(C.c2pa_reader_new())
}

func c2paReaderFromContext(ctx unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.c2pa_reader_from_context((*C.C2paContext)(ctx)))
}

func c2paReaderJson(r unsafe.Pointer) string {
	cs := C.c2pa_reader_json((*C.C2paReader)(r))
	defer C.c2pa_free(unsafe.Pointer(cs))
	return C.GoString(cs)
}

func c2paReaderDetailedJson(r unsafe.Pointer) string {
	cs := C.c2pa_reader_detailed_json((*C.C2paReader)(r))
	defer C.c2pa_free(unsafe.Pointer(cs))
	return C.GoString(cs)
}

func c2paReaderRemoteUrl(r unsafe.Pointer) string {
	cs := C.c2pa_reader_remote_url((*C.C2paReader)(r))
	if cs == nil {
		return ""
	}
	return C.GoString(cs)
}

func c2paReaderIsEmbedded(r unsafe.Pointer) bool {
	return bool(C.c2pa_reader_is_embedded((*C.C2paReader)(r)))
}

func c2paReaderResourceToStream(r unsafe.Pointer, uri string, stream unsafe.Pointer) int64 {
	curi := C.CString(uri)
	defer C.free(unsafe.Pointer(curi))
	return int64(C.c2pa_reader_resource_to_stream((*C.C2paReader)(r), curi, (*C.C2paStream)(stream)))
}

func c2paReaderWithStream(r unsafe.Pointer, format string, stream unsafe.Pointer) unsafe.Pointer {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))
	return unsafe.Pointer(C.c2pa_reader_with_stream((*C.C2paReader)(r), cformat, (*C.C2paStream)(stream)))
}

func c2paReaderWithManifestDataAndStream(r unsafe.Pointer, format string, stream unsafe.Pointer, data []byte) unsafe.Pointer {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))
	var dataPtr *C.uchar
	if len(data) > 0 {
		dataPtr = (*C.uchar)(unsafe.Pointer(&data[0]))
	}
	return unsafe.Pointer(C.c2pa_reader_with_manifest_data_and_stream(
		(*C.C2paReader)(r), cformat, (*C.C2paStream)(stream), dataPtr, C.uintptr_t(len(data))))
}

func c2paReaderWithFragment(r, asset, fragment unsafe.Pointer, format string) unsafe.Pointer {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))
	return unsafe.Pointer(C.c2pa_reader_with_fragment(
		(*C.C2paReader)(r), cformat, (*C.C2paStream)(asset), (*C.C2paStream)(fragment)))
}

// ---- builder ---------------------------------------------------------------

func c2paBuilderFromContext(ctx unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.c2pa_builder_from_context((*C.C2paContext)(ctx)))
}

func c2paBuilderWithDefinition(b unsafe.Pointer, json string) unsafe.Pointer {
	cjson := C.CString(json)
	defer C.free(unsafe.Pointer(cjson))
	return unsafe.Pointer(C.c2pa_builder_with_definition((*C.C2paBuilder)(b), cjson))
}

func c2paBuilderWithArchive(b, stream unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.c2pa_builder_with_archive((*C.C2paBuilder)(b), (*C.C2paStream)(stream)))
}

func c2paBuilderSetNoEmbed(b unsafe.Pointer) {
	C.c2pa_builder_set_no_embed((*C.C2paBuilder)(b))
}

func c2paBuilderSetRemoteUrl(b unsafe.Pointer, url string) int {
	curl := C.CString(url)
	defer C.free(unsafe.Pointer(curl))
	return int(C.c2pa_builder_set_remote_url((*C.C2paBuilder)(b), curl))
}

func c2paBuilderSetBasePath(b unsafe.Pointer, path string) int {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return int(C.c2pa_builder_set_base_path((*C.C2paBuilder)(b), cpath))
}

func c2paBuilderSetIntent(b unsafe.Pointer, intent BuilderIntent, source DigitalSourceType) int {
	ci, ok := intentToC[intent]
	if !ok {
		c2paErrorSetLast("InvalidParameter: unknown builder intent: " + string(intent))
		return -1
	}
	cs, ok := sourceToC[source]
	if !ok {
		c2paErrorSetLast("InvalidParameter: unknown digital source type: " + string(source))
		return -1
	}
	return int(C.c2pa_builder_set_intent((*C.C2paBuilder)(b), uint32(ci), uint32(cs)))
}

func c2paBuilderAddAction(b unsafe.Pointer, json string) int {
	cjson := C.CString(json)
	defer C.free(unsafe.Pointer(cjson))
	return int(C.c2pa_builder_add_action((*C.C2paBuilder)(b), cjson))
}

func c2paBuilderAddResource(b unsafe.Pointer, uri string, stream unsafe.Pointer) int {
	curi := C.CString(uri)
	defer C.free(unsafe.Pointer(curi))
	return int(C.c2pa_builder_add_resource((*C.C2paBuilder)(b), curi, (*C.C2paStream)(stream)))
}

func c2paBuilderAddIngredientFromStream(b unsafe.Pointer, json, format string, stream unsafe.Pointer) int {
	cjson := C.CString(json)
	defer C.free(unsafe.Pointer(cjson))
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))
	return int(C.c2pa_builder_add_ingredient_from_stream((*C.C2paBuilder)(b), cjson, cformat, (*C.C2paStream)(stream)))
}

func c2paBuilderToArchive(b, stream unsafe.Pointer) int {
	return int(C.c2pa_builder_to_archive((*C.C2paBuilder)(b), (*C.C2paStream)(stream)))
}

func c2paBuilderAddIngredientFromArchive(b, stream unsafe.Pointer) int {
	return int(C.c2pa_builder_add_ingredient_from_archive((*C.C2paBuilder)(b), (*C.C2paStream)(stream)))
}

func c2paBuilderWriteIngredientArchive(b unsafe.Pointer, id string, stream unsafe.Pointer) int {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))
	return int(C.c2pa_builder_write_ingredient_archive((*C.C2paBuilder)(b), cid, (*C.C2paStream)(stream)))
}

func c2paBuilderSign(b unsafe.Pointer, format string, input, output, signer unsafe.Pointer) ([]byte, int64) {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))
	var manifest *C.uchar
	n := C.c2pa_builder_sign(
		(*C.C2paBuilder)(b), cformat,
		(*C.C2paStream)(input), (*C.C2paStream)(output),
		(*C.C2paSigner)(signer), &manifest)
	if n < 0 {
		return nil, int64(n)
	}
	return takeCBytes(unsafe.Pointer(manifest), int64(n)), int64(n)
}

func c2paBuilderSignContext(b unsafe.Pointer, format string, input, output unsafe.Pointer) ([]byte, int64) {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))
	var manifest *C.uchar
	n := C.c2pa_builder_sign_context(
		(*C.C2paBuilder)(b), cformat,
		(*C.C2paStream)(input), (*C.C2paStream)(output),
		&manifest)
	if n < 0 {
		return nil, int64(n)
	}
	return takeCBytes(unsafe.Pointer(manifest), int64(n)), int64(n)
}

func c2paBuilderNeedsPlaceholder(b unsafe.Pointer, format string) int {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))
	return int(C.c2pa_builder_needs_placeholder((*C.C2paBuilder)(b), cformat))
}

func c2paBuilderHashType(b unsafe.Pointer, format string) (HashType, int) {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))
	var out uint32
	rc := int(C.c2pa_builder_hash_type((*C.C2paBuilder)(b), cformat, &out))
	if rc < 0 {
		return "", rc
	}
	return hashTypeFromC[C.C2paHashType(out)], rc
}

func c2paBuilderPlaceholder(b unsafe.Pointer, format string) ([]byte, int64) {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))
	var bytesPtr *C.uchar
	n := C.c2pa_builder_placeholder((*C.C2paBuilder)(b), cformat, &bytesPtr)
	if n < 0 {
		return nil, int64(n)
	}
	return takeCBytes(unsafe.Pointer(bytesPtr), int64(n)), int64(n)
}

func c2paBuilderDataHashedPlaceholder(b unsafe.Pointer, reservedSize uintptr, format string) ([]byte, int64) {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))
	var bytesPtr *C.uchar
	n := C.c2pa_builder_data_hashed_placeholder((*C.C2paBuilder)(b), C.uintptr_t(reservedSize), cformat, &bytesPtr)
	if n < 0 {
		return nil, int64(n)
	}
	return takeCBytes(unsafe.Pointer(bytesPtr), int64(n)), int64(n)
}

func c2paBuilderSignDataHashedEmbeddable(b, signer unsafe.Pointer, dataHashJson, format string, asset unsafe.Pointer) ([]byte, int64) {
	cdh := C.CString(dataHashJson)
	defer C.free(unsafe.Pointer(cdh))
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))
	var bytesPtr *C.uchar
	n := C.c2pa_builder_sign_data_hashed_embeddable(
		(*C.C2paBuilder)(b), (*C.C2paSigner)(signer), cdh, cformat,
		(*C.C2paStream)(asset), &bytesPtr)
	if n < 0 {
		return nil, int64(n)
	}
	return takeCBytes(unsafe.Pointer(bytesPtr), int64(n)), int64(n)
}

func c2paBuilderSignEmbeddable(b unsafe.Pointer, format string) ([]byte, int64) {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))
	var bytesPtr *C.uchar
	n := C.c2pa_builder_sign_embeddable((*C.C2paBuilder)(b), cformat, &bytesPtr)
	if n < 0 {
		return nil, int64(n)
	}
	return takeCBytes(unsafe.Pointer(bytesPtr), int64(n)), int64(n)
}

func c2paBuilderSetDataHashExclusions(b unsafe.Pointer, pairs []uint64) int {
	var ptr *C.uint64_t
	if len(pairs) > 0 {
		ptr = (*C.uint64_t)(unsafe.Pointer(&pairs[0]))
	}
	count := len(pairs) / 2
	return int(C.c2pa_builder_set_data_hash_exclusions((*C.C2paBuilder)(b), ptr, C.uintptr_t(count)))
}

func c2paBuilderSetFixedSizeMerkle(b unsafe.Pointer, sizeKB uintptr) int {
	return int(C.c2pa_builder_set_fixed_size_merkle((*C.C2paBuilder)(b), C.uintptr_t(sizeKB)))
}

func c2paBuilderHashMdatBytes(b unsafe.Pointer, mdatId uintptr, data []byte, largeSize bool) int {
	var dataPtr *C.uchar
	if len(data) > 0 {
		dataPtr = (*C.uchar)(unsafe.Pointer(&data[0]))
	}
	return int(C.c2pa_builder_hash_mdat_bytes(
		(*C.C2paBuilder)(b), C.uintptr_t(mdatId), dataPtr, C.uintptr_t(len(data)), C.bool(largeSize)))
}

func c2paBuilderUpdateHashFromStream(b unsafe.Pointer, format string, stream unsafe.Pointer) int {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))
	return int(C.c2pa_builder_update_hash_from_stream((*C.C2paBuilder)(b), cformat, (*C.C2paStream)(stream)))
}

func c2paFormatEmbeddable(format string, manifestBytes []byte) ([]byte, int64) {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))
	var inPtr *C.uchar
	if len(manifestBytes) > 0 {
		inPtr = (*C.uchar)(unsafe.Pointer(&manifestBytes[0]))
	}
	var outPtr *C.uchar
	n := C.c2pa_format_embeddable(cformat, inPtr, C.uintptr_t(len(manifestBytes)), &outPtr)
	if n < 0 {
		return nil, int64(n)
	}
	return takeCBytes(unsafe.Pointer(outPtr), int64(n)), int64(n)
}

// ---- signer ----------------------------------------------------------------

func c2paSignerReserveSize(signer unsafe.Pointer) int64 {
	return int64(C.c2pa_signer_reserve_size((*C.C2paSigner)(signer)))
}

func c2paSignerCreate(handle uintptr, alg SigningAlg, tsaUrl, certificates string) unsafe.Pointer {
	cAlg, ok := signingAlgToC[alg]
	if !ok {
		c2paErrorSetLast("InvalidParameter: unknown signing algorithm: " + string(alg))
		return nil
	}
	cTsa := C.CString(tsaUrl)
	defer C.free(unsafe.Pointer(cTsa))
	cCerts := C.CString(certificates)
	defer C.free(unsafe.Pointer(cCerts))
	return unsafe.Pointer(C.create_signer(C.uintptr_t(handle), cAlg, cTsa, cCerts))
}

func c2paSignerFromInfo(alg SigningAlg, signCert, privateKey, timestampUrl string) unsafe.Pointer {
	cAlg := C.CString(string(alg))
	defer C.free(unsafe.Pointer(cAlg))
	cCert := C.CString(signCert)
	defer C.free(unsafe.Pointer(cCert))
	cKey := C.CString(privateKey)
	defer C.free(unsafe.Pointer(cKey))
	var cTaUrl *C.char
	if timestampUrl != "" {
		cTaUrl = C.CString(timestampUrl)
		defer C.free(unsafe.Pointer(cTaUrl))
	}
	cInfo := C.C2paSignerInfo{
		cAlg,
		cCert,
		cKey,
		cTaUrl,
	}
	return unsafe.Pointer(C.c2pa_signer_from_info(&cInfo))
}

func c2paIdentitySignerCreate(c2paSigner, identitySigner unsafe.Pointer, referencedAssertions, roles []string) unsafe.Pointer {
	crefs, freeRefs := cStringArray(referencedAssertions)
	defer freeRefs()
	rolesPtr, freeRoles := cStringArray(roles)
	defer freeRoles()
	return unsafe.Pointer(C.c2pa_identity_signer_create(
		(*C.C2paSigner)(c2paSigner), (*C.C2paSigner)(identitySigner), crefs, rolesPtr))
}

func c2paEd25519Sign(input []byte, privateKey string) []byte {
	cKey := C.CString(privateKey)
	defer C.free(unsafe.Pointer(cKey))
	var inPtr *C.uchar
	if len(input) > 0 {
		inPtr = (*C.uchar)(unsafe.Pointer(&input[0]))
	}
	sig := C.c2pa_ed25519_sign(inPtr, C.uintptr_t(len(input)), cKey)
	if sig == nil {
		return nil
	}
	defer C.c2pa_free(unsafe.Pointer(sig))
	return C.GoBytes(unsafe.Pointer(sig), C.int(ed25519SignatureLen))
}

func cStringArray(values []string) (**C.char, func()) {
	if len(values) == 0 {
		return nil, func() {}
	}
	ptrs := make([]*C.char, len(values)+1)
	for i, value := range values {
		ptrs[i] = C.CString(value)
	}
	return &ptrs[0], func() {
		for _, ptr := range ptrs[:len(values)] {
			C.free(unsafe.Pointer(ptr))
		}
	}
}

// ---- http resolver ---------------------------------------------------------

func c2paHttpResolverCreate(handle uintptr) unsafe.Pointer {
	return unsafe.Pointer(C.create_http_resolver(C.uintptr_t(handle)))
}

// ---- settings --------------------------------------------------------------

func c2paSettingsNew() unsafe.Pointer {
	return unsafe.Pointer(C.c2pa_settings_new())
}

func c2paSettingsUpdateFromString(s unsafe.Pointer, content, format string) int {
	cContent := C.CString(content)
	defer C.free(unsafe.Pointer(cContent))
	cFormat := C.CString(format)
	defer C.free(unsafe.Pointer(cFormat))
	return int(C.c2pa_settings_update_from_string((*C.C2paSettings)(s), cContent, cFormat))
}

func c2paSettingsSetValue(s unsafe.Pointer, path, value string) int {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	return int(C.c2pa_settings_set_value((*C.C2paSettings)(s), cPath, cValue))
}

// ---- stream ----------------------------------------------------------------

func c2paCreateStream(handle uintptr) unsafe.Pointer {
	return unsafe.Pointer(C.create_stream(C.uintptr_t(handle)))
}

func c2paReleaseStream(s unsafe.Pointer) {
	C.c2pa_release_stream((*C.C2paStream)(s))
}

// ---- exported callbacks ----------------------------------------------------

//export signerCallback
func signerCallback(context C.uintptr_t, input *C.uint8_t, inputSize C.uintptr_t, output *C.uint8_t, outputSize C.uintptr_t) C.intptr_t {
	in := unsafe.Slice((*byte)(unsafe.Pointer(input)), int(inputSize))
	out := unsafe.Slice((*byte)(unsafe.Pointer(output)), int(outputSize))
	n, ok := goSignerCallback(uintptr(context), in, out)
	if !ok {
		return C.intptr_t(-1)
	}
	return C.intptr_t(n)
}

//export streamRead
func streamRead(context C.uintptr_t, buffer *C.uint8_t, size C.intptr_t) C.intptr_t {
	slice := unsafe.Slice((*byte)(unsafe.Pointer(buffer)), int(size))
	return C.intptr_t(goStreamRead(uintptr(context), slice))
}

//export streamSeek
func streamSeek(context C.uintptr_t, offset C.intptr_t, mode C.C2paSeekMode) C.intptr_t {
	return C.intptr_t(goStreamSeek(uintptr(context), int64(offset), int(mode)))
}

//export streamWrite
func streamWrite(context C.uintptr_t, buffer *C.uint8_t, size C.intptr_t) C.intptr_t {
	slice := unsafe.Slice((*byte)(unsafe.Pointer(buffer)), int(size))
	return C.intptr_t(goStreamWrite(uintptr(context), slice))
}

//export streamFlush
func streamFlush(context C.uintptr_t) C.intptr_t {
	return C.intptr_t(goStreamFlush(uintptr(context)))
}

//export progressCallback
func progressCallback(context C.uintptr_t, phase C.C2paProgressPhase, step C.uint32_t, total C.uint32_t) C.int {
	phaseStr, ok := progressPhaseFromC[phase]
	if !ok {
		return C.int(0)
	}
	if goProgressCallback(uintptr(context), phaseStr, uint32(step), uint32(total)) {
		return C.int(1)
	}
	return C.int(0)
}

//export httpResolverCallback
func httpResolverCallback(context C.uintptr_t, request *C.C2paHttpRequest, response *C.C2paHttpResponse) C.int {
	if request == nil || response == nil {
		c2paErrorSetLast("InvalidParameter: nil http request or response")
		return -1
	}
	url := C.GoString(request.url)
	method := C.GoString(request.method)
	headers := ""
	if request.headers != nil {
		headers = C.GoString(request.headers)
	}
	var body []byte
	if request.body != nil && request.body_len > 0 {
		body = C.GoBytes(unsafe.Pointer(request.body), C.int(request.body_len))
	}
	status, respBody, errMsg := goHttpResolve(uintptr(context), url, method, headers, body)
	if errMsg != "" {
		c2paErrorSetLast(errMsg)
		return -1
	}
	response.status = C.int32_t(status)
	response.body_len = C.uintptr_t(len(respBody))
	if len(respBody) > 0 {
		buf := C.malloc(C.size_t(len(respBody)))
		if buf == nil {
			c2paErrorSetLast("Other: malloc failed for http response body")
			return -1
		}
		C.memcpy(buf, unsafe.Pointer(&respBody[0]), C.size_t(len(respBody)))
		response.body = (*C.uchar)(buf)
	} else {
		response.body = nil
	}
	return 0
}
