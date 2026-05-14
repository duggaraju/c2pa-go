package c2pa

// #include <stdlib.h>
// #include <string.h>
// #include "c2pa_helper.h"
import "C"

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"runtime/cgo"
	"strings"
	"unsafe"
)

// HttpResolver is the Go-side counterpart of the C2paHttpResolver. It receives
// HTTP requests issued by the C2PA SDK (e.g. for remote-manifest fetches,
// OCSP, or RFC 3161 timestamp authorities) and is expected to perform them
// using Go's network stack. This removes the need to compile the libc2pa_c
// library with the reqwest-based `http` feature.
type HttpResolver interface {
	// Resolve performs the given request and returns the response. The
	// caller retains ownership of req; implementations must not mutate it.
	Resolve(req *http.Request) (*http.Response, error)
}

// HttpResolverAdapter bridges a Go HttpResolver to a C2paHttpResolver*. The
// adapter is consumed by ContextBuilder.SetHttpResolver; after a successful
// call to SetHttpResolver the adapter's C pointer is owned by the builder and
// must not be freed by the caller. The Go-side cgo.Handle is released when
// Close is called.
type HttpResolverAdapter struct {
	resolver HttpResolver
	ptr      *C.C2paHttpResolver
	handle   cgo.Handle
}

// DefaultHttpResolver is a simple HttpResolver backed by an *http.Client.
// If Client is nil, http.DefaultClient is used.
type DefaultHttpResolver struct {
	Client *http.Client
}

// Resolve implements HttpResolver using the configured *http.Client.
func (d *DefaultHttpResolver) Resolve(req *http.Request) (*http.Response, error) {
	client := d.Client
	if client == nil {
		client = http.DefaultClient
	}
	return client.Do(req)
}

// NewHttpResolver wraps a Go HttpResolver in a C-callable resolver. The
// returned adapter must be either passed to ContextBuilder.SetHttpResolver
// (which takes ownership of the C pointer) or released with Close.
func NewHttpResolver(resolver HttpResolver) (*HttpResolverAdapter, error) {
	if resolver == nil {
		return nil, fmt.Errorf("resolver is nil")
	}
	a := &HttpResolverAdapter{resolver: resolver}
	a.handle = cgo.NewHandle(a)
	a.ptr = C.create_http_resolver(C.uintptr_t(a.handle))
	if a.ptr == nil {
		a.handle.Delete()
		return nil, fmt.Errorf("failed to create c2pa http resolver: %s", c2paError())
	}
	return a, nil
}

// Close releases the underlying C resolver if it has not been consumed by
// ContextBuilder.SetHttpResolver, and always releases the Go-side handle.
// Safe to call multiple times.
func (a *HttpResolverAdapter) Close() {
	if a == nil {
		return
	}
	if a.ptr != nil {
		C.c2pa_free(unsafe.Pointer(a.ptr))
		a.ptr = nil
	}
	if a.handle != 0 {
		a.handle.Delete()
		a.handle = 0
	}
}

// parseHeaders parses the "Name: Value\n" formatted header block produced by
// the Rust side into a textproto-style map suitable for http.Request.Header.
func parseHeaders(raw string) http.Header {
	h := http.Header{}
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}
		idx := strings.IndexByte(line, ':')
		if idx < 0 {
			continue
		}
		name := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		if name == "" {
			continue
		}
		h.Add(name, value)
	}
	return h
}

//export httpResolverCallback
func httpResolverCallback(context C.uintptr_t, request *C.C2paHttpRequest, response *C.C2paHttpResponse) C.int {
	if request == nil || response == nil {
		setLastError("InvalidParameter: nil http request or response")
		return -1
	}
	adapter, ok := cgo.Handle(context).Value().(*HttpResolverAdapter)
	if !ok || adapter == nil || adapter.resolver == nil {
		setLastError("Other: http resolver handle is invalid")
		return -1
	}

	url := C.GoString(request.url)
	method := C.GoString(request.method)
	headers := ""
	if request.headers != nil {
		headers = C.GoString(request.headers)
	}
	var bodyReader io.Reader
	if request.body != nil && request.body_len > 0 {
		body := C.GoBytes(unsafe.Pointer(request.body), C.int(request.body_len))
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		setLastError(fmt.Sprintf("Other: %s", err))
		return -1
	}
	req.Header = parseHeaders(headers)

	resp, err := adapter.resolver.Resolve(req)
	if err != nil {
		setLastError(fmt.Sprintf("Other: %s", err))
		return -1
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		setLastError(fmt.Sprintf("Other: %s", err))
		return -1
	}

	response.status = C.int32_t(resp.StatusCode)
	response.body_len = C.uintptr_t(len(body))
	if len(body) > 0 {
		// The Rust side takes ownership of `body` and frees it with
		// libc::free, so we must allocate via the C allocator.
		buf := C.malloc(C.size_t(len(body)))
		if buf == nil {
			setLastError("Other: malloc failed for http response body")
			return -1
		}
		C.memcpy(buf, unsafe.Pointer(&body[0]), C.size_t(len(body)))
		response.body = (*C.uchar)(buf)
	} else {
		response.body = nil
	}
	return 0
}

// setLastError forwards a Go-side error message to the Rust side via
// c2pa_error_set_last so that it surfaces through C2paError().
func setLastError(msg string) {
	cs := C.CString(msg)
	defer C.free(unsafe.Pointer(cs))
	C.c2pa_error_set_last(cs)
}
