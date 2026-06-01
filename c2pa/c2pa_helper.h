#ifndef C2PA_HELPER_H
#define C2PA_HELPER_H

#include <stdint.h>
#include <c2pa.h>

#ifdef __cplusplus
extern "C" {
#endif

#ifdef __GNUC__
#define WEAK __attribute__((weak))
#else
#define WEAK
#endif

#ifdef _WIN32
#define CGO_EXPORT __declspec(dllexport)
#else
#define CGO_EXPORT
#endif

extern CGO_EXPORT intptr_t signerCallback(uintptr_t context, uint8_t* input, uintptr_t input_size, uint8_t* output, uintptr_t output_size);

WEAK intptr_t signer_callback(const void* context, const uint8_t* input, uintptr_t input_size, uint8_t* output, uintptr_t output_size)
{
	return signerCallback((uintptr_t)context, (uint8_t*)input, input_size, output, output_size);
}

WEAK C2paSigner* create_signer(uintptr_t context, C2paSigningAlg alg, const char* tsa_url, const char* certificates)
{
    return c2pa_signer_create((const void*)context, (SignerCallback)signer_callback, alg, certificates, tsa_url);
}

extern CGO_EXPORT intptr_t streamRead(uintptr_t context, uint8_t *buffer, intptr_t size);
WEAK intptr_t stream_read(StreamContext *context, uint8_t *buffer, intptr_t size) {

	return streamRead((uintptr_t)context, buffer, size);
}

extern CGO_EXPORT intptr_t streamWrite(uintptr_t context, uint8_t* buffer, intptr_t size);
WEAK intptr_t stream_write(StreamContext *context, const uint8_t *buffer, intptr_t size) {
	return streamWrite((uintptr_t)context, (uint8_t*)buffer, size);
}

extern CGO_EXPORT intptr_t streamSeek(uintptr_t context, intptr_t offset, C2paSeekMode mode);
WEAK intptr_t stream_seek(StreamContext *context, intptr_t offset, C2paSeekMode mode) {
	return streamSeek((uintptr_t)context, offset, mode);
}

extern CGO_EXPORT intptr_t streamFlush(uintptr_t context);
WEAK intptr_t stream_flush(StreamContext *context) {
	return streamFlush((uintptr_t)context);
}

WEAK C2paStream* create_stream(uintptr_t context) {
	return c2pa_create_stream((StreamContext*)context, stream_read, stream_seek, stream_write, stream_flush);
}

extern CGO_EXPORT int httpResolverCallback(uintptr_t context, C2paHttpRequest* request, C2paHttpResponse* response);

WEAK int http_resolver_callback(void* context, const C2paHttpRequest* request, C2paHttpResponse* response) {
	return httpResolverCallback((uintptr_t)context, (C2paHttpRequest*)request, response);
}

WEAK C2paHttpResolver* create_http_resolver(uintptr_t context) {
	return c2pa_http_resolver_create((const void*)context, (C2paHttpResolverCallback)http_resolver_callback);
}

extern CGO_EXPORT int progressCallback(uintptr_t context, C2paProgressPhase phase, uint32_t step, uint32_t total);

WEAK int progress_callback(const void* context, C2paProgressPhase phase, uint32_t step, uint32_t total) {
	return progressCallback((uintptr_t)context, phase, step, total);
}

WEAK int set_progress_callback(C2paContextBuilder* builder, uintptr_t context) {
	return c2pa_context_builder_set_progress_callback(builder, (const void*)context, (ProgressCCallback)progress_callback);
}

#ifdef __cplusplus
}
#endif

#endif