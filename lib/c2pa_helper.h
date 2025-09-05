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

extern intptr_t GoSignerCallback(uintptr_t context, uint8_t* input, uintptr_t input_size, uint8_t* output, uintptr_t output_size);

WEAK intptr_t signer_callback(const void* context, const uint8_t* input, uintptr_t input_size, uint8_t* output, uintptr_t output_size)
{
	return GoSignerCallback((uintptr_t)context, (uint8_t*)input, input_size, output, output_size);
}

WEAK C2paSigner* create_signer(uintptr_t context, C2paSigningAlg alg, const char* tsa_url, const char* certificates)
{
    return c2pa_signer_create((const void*)context, (SignerCallback)signer_callback, alg, tsa_url, certificates);
}

WEAK intptr_t sign_data(C2paBuilder* builder, const char* format, C2paStream* input, C2paStream* output, C2paSigner* signer, void* manifest)
{
	return c2pa_builder_sign(builder, format, input, output, signer, (const uint8_t**)manifest);
}


extern intptr_t StreamRead(uintptr_t context, uint8_t *buffer, intptr_t size);
WEAK intptr_t stream_read(StreamContext *context, uint8_t *buffer, intptr_t size) {

	return StreamRead((uintptr_t)context, buffer, size);
}

extern intptr_t StreamWrite(uintptr_t context, uint8_t* buffer, intptr_t size);
WEAK intptr_t stream_write(StreamContext *context, const uint8_t *buffer, intptr_t size) {
	return StreamWrite((uintptr_t)context, (uint8_t*)buffer, size);
}

extern intptr_t StreamSeek(uintptr_t context, intptr_t offset, C2paSeekMode mode);
WEAK intptr_t stream_seek(StreamContext *context, intptr_t offset, C2paSeekMode mode) {
	return StreamSeek((uintptr_t)context, offset, mode);
}

extern intptr_t StreamFlush(uintptr_t context);
WEAK intptr_t stream_flush(StreamContext *context) {
	return StreamFlush((uintptr_t)context);
}

WEAK C2paStream* create_stream(uintptr_t context) {
	return c2pa_create_stream((StreamContext*)context, stream_read, stream_seek, stream_write, stream_flush);
}

#ifdef __cplusplus
}
#endif

#endif