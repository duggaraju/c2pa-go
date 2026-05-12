package lib

// #include <c2pa.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// Context wraps an immutable C2paContext*. A Context is shareable and may be
// used to create multiple Reader and Builder instances.
type Context struct {
	ptr *C.C2paContext
}

// Ptr returns the underlying C pointer. Use carefully; callers must not keep
// the pointer past the lifetime of the Context.
func (c *Context) Ptr() *C.C2paContext {
	return c.ptr
}

// NewContext creates a new immutable Context with default settings.
func NewContext() (*Context, error) {
	ptr := C.c2pa_context_new()
	if ptr == nil {
		return nil, fmt.Errorf("failed to create c2pa context: %s", C2paError())
	}
	return &Context{ptr: ptr}, nil
}

// Close releases the underlying C context. Safe to call once; subsequent
// calls are no-ops.
func (c *Context) Close() {
	if c.ptr != nil {
		C.c2pa_free(unsafe.Pointer(c.ptr))
		c.ptr = nil
	}
}

// Cancel requests cancellation of any in-progress operation on this context.
// Thread-safe.
func (c *Context) Cancel() error {
	if rc := C.c2pa_context_cancel(c.ptr); rc != 0 {
		return fmt.Errorf("failed to cancel c2pa context: %s", C2paError())
	}
	return nil
}
