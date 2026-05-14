package c2pa

// #include <c2pa.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// ContextBuilder wraps a C2paContextBuilder*. It is consumed by Build().
type ContextBuilder struct {
	ptr *C.C2paContextBuilder
}

// NewContextBuilder creates a new context builder with default settings.
func NewContextBuilder() (*ContextBuilder, error) {
	ptr := C.c2pa_context_builder_new()
	if ptr == nil {
		return nil, fmt.Errorf("failed to create c2pa context builder: %s", c2paError())
	}
	return &ContextBuilder{ptr: ptr}, nil
}

// Close releases the underlying builder if it has not been consumed by Build.
func (b *ContextBuilder) Close() {
	if b.ptr != nil {
		C.c2pa_free(unsafe.Pointer(b.ptr))
		b.ptr = nil
	}
}

// SetSigner attaches a signer to the context. The builder wraps the supplied
// Signer in an internal C-callable adapter and takes ownership of it; the
// adapter is kept alive for the lifetime of the Context produced by Build()
// and released when that Context is closed.
func (b *ContextBuilder) SetSigner(signer Signer) error {
	if b.ptr == nil {
		return fmt.Errorf("context builder is closed")
	}
	if signer == nil {
		return fmt.Errorf("signer is nil")
	}
	adapter, err := newSigner(signer)
	if err != nil {
		return err
	}
	if rc := C.c2pa_context_builder_set_signer(b.ptr, adapter.ptr); rc != 0 {
		adapter.Close()
		return fmt.Errorf("failed to set signer on context builder: %s", c2paError())
	}
	// The builder consumed the C signer pointer; prevent double free.
	adapter.ptr = nil
	return nil
}

// SetSignerInfo creates a C2PA signer from a local certificate + private key
// (and optional RFC 3161 timestamp URL) and attaches it to the context.
func (b *ContextBuilder) SetSignerInfo(info SignerInfo) error {
	if b.ptr == nil {
		return fmt.Errorf("context builder is closed")
	}
	cAlg := C.CString(info.Alg)
	defer C.free(unsafe.Pointer(cAlg))
	cCert := C.CString(info.SignCert)
	defer C.free(unsafe.Pointer(cCert))
	cKey := C.CString(info.PrivateKey)
	defer C.free(unsafe.Pointer(cKey))
	var cTaUrl *C.char
	if info.TimestampUrl != "" {
		cTaUrl = C.CString(info.TimestampUrl)
		defer C.free(unsafe.Pointer(cTaUrl))
	}

	cInfo := C.C2paSignerInfo{
		alg:         cAlg,
		sign_cert:   cCert,
		private_key: cKey,
		ta_url:      cTaUrl,
	}
	cSigner := C.c2pa_signer_from_info(&cInfo)
	if cSigner == nil {
		return fmt.Errorf("failed to create signer from info: %s", c2paError())
	}
	if rc := C.c2pa_context_builder_set_signer(b.ptr, cSigner); rc != 0 {
		C.c2pa_free(unsafe.Pointer(cSigner))
		return fmt.Errorf("failed to set signer on context builder: %s", c2paError())
	}
	return nil
}

// SetSettings configures the builder with the given settings. Settings are
// cloned internally; the caller retains ownership and should still Close() it.
func (b *ContextBuilder) SetSettings(settings *Settings) error {
	if b.ptr == nil {
		return fmt.Errorf("context builder is closed")
	}
	if settings == nil || settings.ptr == nil {
		return fmt.Errorf("settings is nil")
	}
	if rc := C.c2pa_context_builder_set_settings(b.ptr, settings.ptr); rc != 0 {
		return fmt.Errorf("failed to set settings on context builder: %s", c2paError())
	}
	return nil
}

// SetHttpResolver attaches a custom HTTP resolver to the context. The builder
// takes ownership of the underlying C resolver; the adapter's C pointer is
// cleared so its Close() will not double-free it. The Go-side handle is still
// released when the adapter is closed (or when the owning Context is freed,
// if the caller keeps a reference to the adapter alive that long).
//
// Because the resolver is consumed on success, callers should ensure that
// the HttpResolverAdapter (and thus its cgo.Handle) outlives any Context
// built from this builder — typically by keeping a reference to the adapter
// next to the Context and calling Close on both at teardown.
func (b *ContextBuilder) SetHttpResolver(resolver *HttpResolverAdapter) error {
	if b.ptr == nil {
		return fmt.Errorf("context builder is closed")
	}
	if resolver == nil || resolver.ptr == nil {
		return fmt.Errorf("http resolver is nil")
	}
	if rc := C.c2pa_context_builder_set_http_resolver(b.ptr, resolver.ptr); rc != 0 {
		return fmt.Errorf("failed to set http resolver on context builder: %s", c2paError())
	}
	// The builder consumed the C resolver pointer; prevent double free.
	resolver.ptr = nil
	return nil
}

// Build consumes the builder and returns an immutable Context.
func (b *ContextBuilder) Build() (*Context, error) {
	if b.ptr == nil {
		return nil, fmt.Errorf("context builder is closed")
	}
	ptr := C.c2pa_context_builder_build(b.ptr)
	// The builder is consumed regardless of success.
	b.ptr = nil
	if ptr == nil {
		return nil, fmt.Errorf("failed to build c2pa context: %s", c2paError())
	}
	return &Context{ptr: ptr}, nil
}
