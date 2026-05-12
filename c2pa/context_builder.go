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
		return nil, fmt.Errorf("failed to create c2pa context builder: %s", C2paError())
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

// SetSigner attaches a signer to the context. The builder takes ownership of
// the underlying C signer; the SignerAdapter's pointer is cleared so its
// Close() will not double-free it. The Go-side handle is still released when
// the SignerAdapter is closed.
func (b *ContextBuilder) SetSigner(signer *SignerAdapter) error {
	if b.ptr == nil {
		return fmt.Errorf("context builder is closed")
	}
	if signer == nil || signer.ptr == nil {
		return fmt.Errorf("signer is nil")
	}
	if rc := C.c2pa_context_builder_set_signer(b.ptr, signer.ptr); rc != 0 {
		return fmt.Errorf("failed to set signer on context builder: %s", C2paError())
	}
	// The builder consumed the C signer pointer; prevent double free.
	signer.ptr = nil
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
		return fmt.Errorf("failed to set settings on context builder: %s", C2paError())
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
		return fmt.Errorf("failed to set http resolver on context builder: %s", C2paError())
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
		return nil, fmt.Errorf("failed to build c2pa context: %s", C2paError())
	}
	return &Context{ptr: ptr}, nil
}
