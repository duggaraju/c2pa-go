package c2pa

import (
	"fmt"
	"runtime/cgo"
	"unsafe"
)

// ContextBuilder wraps a C2paContextBuilder*. It is consumed by Build().
type ContextBuilder struct {
	ptr      unsafe.Pointer
	cleanups []func()
}

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

// goProgressCallback is invoked from the cgo progressCallback in native.go.
func goProgressCallback(handle uintptr, phase ProgressPhase, step, total uint32) bool {
	callback, ok := cgo.Handle(handle).Value().(ProgressCallback)
	if !ok || callback == nil {
		return false
	}
	return callback.Progress(phase, step, total)
}

// NewContextBuilder creates a new context builder with default settings.
func NewContextBuilder() (*ContextBuilder, error) {
	ptr := c2paContextBuilderNew()
	if ptr == nil {
		return nil, fmt.Errorf("failed to create c2pa context builder: %s", c2paError())
	}
	return &ContextBuilder{ptr: ptr}, nil
}

// Close releases the underlying builder if it has not been consumed by Build.
func (b *ContextBuilder) Close() {
	if b.ptr != nil {
		c2paFree(b.ptr)
		b.ptr = nil
	}
	for _, cleanup := range b.cleanups {
		cleanup()
	}
	b.cleanups = nil
}

// SetSigner attaches a signer to the context. Native-backed signers transfer
// ownership of their native signer to the builder; callback signers are
// wrapped in an internal C-callable adapter. Any callback handles are kept
// alive for the lifetime of the Context produced by Build().
func (b *ContextBuilder) SetSigner(signer Signer) error {
	if b.ptr == nil {
		return fmt.Errorf("context builder is closed")
	}
	if signer == nil {
		return fmt.Errorf("signer is nil")
	}
	native, err := takeNativeSigner(signer)
	if err != nil {
		return err
	}
	if rc := c2paContextBuilderSetSigner(b.ptr, native.ptr); rc != 0 {
		native.Close()
		return fmt.Errorf("failed to set signer on context builder: %s", c2paError())
	}
	b.takeHandles(native.handles)
	native.handles = nil
	// The builder consumed the C signer pointer; prevent double free.
	native.ptr = nil
	return nil
}

// SetSignerInfo creates a C2PA signer from a local certificate + private key
// (and optional RFC 3161 timestamp URL) and attaches it to the context.
func (b *ContextBuilder) SetSignerInfo(info SignerInfo) error {
	if b.ptr == nil {
		return fmt.Errorf("context builder is closed")
	}
	signer, err := NewSignerFromInfo(info)
	if err != nil {
		return fmt.Errorf("failed to create signer from info: %w", err)
	}
	if err := b.SetSigner(signer); err != nil {
		return err
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
	if rc := c2paContextBuilderSetSettings(b.ptr, settings.ptr); rc != 0 {
		return fmt.Errorf("failed to set settings on context builder: %s", c2paError())
	}
	return nil
}

// SetHttpResolver attaches a custom HTTP resolver to the context. The builder
// wraps the supplied resolver in an internal C-callable adapter, takes
// ownership of the underlying C resolver pointer, and keeps the Go-side
// handle alive for the lifetime of the Context produced by Build().
func (b *ContextBuilder) SetHttpResolver(resolver HttpResolver) error {
	if b.ptr == nil {
		return fmt.Errorf("context builder is closed")
	}
	if resolver == nil {
		return fmt.Errorf("http resolver is nil")
	}
	adapter, err := newHttpResolver(resolver)
	if err != nil {
		return err
	}
	if rc := c2paContextBuilderSetHttpResolver(b.ptr, adapter.ptr); rc != 0 {
		adapter.Close()
		return fmt.Errorf("failed to set http resolver on context builder: %s", c2paError())
	}
	b.takeHandle(adapter.handle)
	adapter.handle = 0
	// The builder consumed the C resolver pointer; prevent double free.
	adapter.ptr = nil
	return nil
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
	if rc := c2paContextBuilderSetProgressCallback(b.ptr, uintptr(handle)); rc != 0 {
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
	ptr := c2paContextBuilderBuild(b.ptr)
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
