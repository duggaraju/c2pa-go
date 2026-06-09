package c2pa

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
	ptr      unsafe.Pointer
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
	a.ptr = c2paHttpResolverCreate(uintptr(a.handle))
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
		c2paFree(a.ptr)
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

// goHttpResolve is invoked from the cgo httpResolverCallback in native.go.
// Returns (status, body, errMsg); errMsg is non-empty on failure.
func goHttpResolve(handle uintptr, url, method, headers string, body []byte) (int, []byte, string) {
	adapter, ok := cgo.Handle(handle).Value().(*HttpResolverAdapter)
	if !ok || adapter == nil || adapter.resolver == nil {
		return 0, nil, "Other: http resolver handle is invalid"
	}

	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return 0, nil, fmt.Sprintf("Other: %s", err)
	}
	req.Header = parseHeaders(headers)

	resp, err := adapter.resolver.Resolve(req)
	if err != nil {
		return 0, nil, fmt.Sprintf("Other: %s", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Sprintf("Other: %s", err)
	}
	return resp.StatusCode, respBody, ""
}
