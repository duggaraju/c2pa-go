package c2pa

// #include "c2pa_helper.h"
import "C"

import (
	"fmt"
	"runtime/cgo"
	"unsafe"
)

// ContextBuilder wraps a C2paContextBuilder*. It is consumed by Build().
type ContextBuilder struct {
	ptr      *C.C2paContextBuilder
	cleanups []func()
}

// ProgressPhase mirrors C2paProgressPhase and identifies the operation stage
// reported by a progress callback.
type ProgressPhase C.C2paProgressPhase

const (
	ProgressReading                ProgressPhase = C.Reading
	ProgressVerifyingManifest      ProgressPhase = C.VerifyingManifest
	ProgressVerifyingSignature     ProgressPhase = C.VerifyingSignature
	ProgressVerifyingIngredient    ProgressPhase = C.VerifyingIngredient
	ProgressVerifyingAssetHash     ProgressPhase = C.VerifyingAssetHash
	ProgressAddingIngredient       ProgressPhase = C.AddingIngredient
	ProgressThumbnail              ProgressPhase = C.Thumbnail
	ProgressHashing                ProgressPhase = C.Hashing
	ProgressSigning                ProgressPhase = C.Signing
	ProgressEmbedding              ProgressPhase = C.Embedding
	ProgressFetchingRemoteManifest ProgressPhase = C.FetchingRemoteManifest
	ProgressWriting                ProgressPhase = C.Writing
	ProgressFetchingOCSP           ProgressPhase = C.FetchingOCSP
	ProgressFetchingTimestamp      ProgressPhase = C.FetchingTimestamp
)

// ProgressCallback receives progress updates from the native SDK. Return true
// to continue or false to request cancellation.
type ProgressCallback interface {
	Progress(phase ProgressPhase, step uint32, total uint32) bool
}

// ProgressFunc adapts a function to the ProgressCallback interface.
type ProgressFunc func(phase ProgressPhase, step uint32, total uint32) bool

// Progress implements ProgressCallback.
func (f ProgressFunc) Progress(phase ProgressPhase, step uint32, total uint32) bool {
	if f == nil {
		return true
	}
	return f(phase, step, total)
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
	for _, cleanup := range b.cleanups {
		cleanup()
	}
	b.cleanups = nil
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
	b.takeHandles(adapter.handles)
	adapter.handles = nil
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
	cSigner, err := NewSignerFromInfo(info)
	if err != nil {
		return fmt.Errorf("failed to create signer from info: %s", c2paError())
	}

	if rc := C.c2pa_context_builder_set_signer(b.ptr, cSigner.ptr); rc != 0 {
		cSigner.Close()
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
	b.takeHandle(resolver.handle)
	resolver.handle = 0
	// The builder consumed the C resolver pointer; prevent double free.
	resolver.ptr = nil
	return nil
}

//export progressCallback
func progressCallback(context C.uintptr_t, phase C.C2paProgressPhase, step C.uint32_t, total C.uint32_t) C.int {
	handle := cgo.Handle(context)
	callback, ok := handle.Value().(ProgressCallback)
	if !ok || callback == nil {
		return C.int(0)
	}
	if callback.Progress(ProgressPhase(phase), uint32(step), uint32(total)) {
		return C.int(1)
	}
	return C.int(0)
}

// SetProgressCallback attaches a Go progress callback to the context builder.
// The callback remains active for the lifetime of the Context built from this
// builder, or until the builder is closed without building.
func (b *ContextBuilder) SetProgressCallback(callback ProgressCallback) error {
	if b.ptr == nil {
		return fmt.Errorf("context builder is closed")
	}
	if callback == nil {
		return fmt.Errorf("progress callback is nil")
	}
	handle := cgo.NewHandle(callback)
	if rc := C.set_progress_callback(b.ptr, C.uintptr_t(handle)); rc != 0 {
		handle.Delete()
		return fmt.Errorf("failed to set progress callback on context builder: %s", c2paError())
	}
	b.takeHandle(handle)
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
		for _, cleanup := range b.cleanups {
			cleanup()
		}
		b.cleanups = nil
		return nil, fmt.Errorf("failed to build c2pa context: %s", c2paError())
	}
	ctx := &Context{ptr: ptr, cleanups: b.cleanups}
	b.cleanups = nil
	return ctx, nil
}

func (b *ContextBuilder) takeHandle(handle cgo.Handle) {
	if handle == 0 {
		return
	}
	b.cleanups = append(b.cleanups, func() {
		handle.Delete()
	})
}

func (b *ContextBuilder) takeHandles(handles []cgo.Handle) {
	for _, handle := range handles {
		b.takeHandle(handle)
	}
}
