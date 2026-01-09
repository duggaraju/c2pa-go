//go:build linux && cgo
// +build linux,cgo

package lib

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCpaVersion(t *testing.T) {
	v := CpaVersion()
	assert.NotEmpty(t, v)
	assert.Regexp(t, regexp.MustCompile(`^c2pa-c-ffi/\d+\.\d+\.\d+\s+c2pa-rs/\d+\.\d+\.\d+$`), v)
}

func TestC2paError_AfterReaderFailure(t *testing.T) {
	// Trigger a failure that happens inside the C2PA library (so c2pa_error is populated),
	// not a Go-side error like "os.Open" failing.
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.jpg")
	err := os.WriteFile(path, nil, 0o644)
	assert.NoError(t, err)

	_, readerErr := ReaderFromFile(path)
	assert.Error(t, readerErr)

	lastErr := C2paError()
	assert.NotEmpty(t, lastErr)
}

func TestReaderFromFile_NotFound(t *testing.T) {
	_, err := ReaderFromFile("/nonexistent/file/path.jpg")
	if err == nil {
		t.Error("ReaderFromFile should fail for nonexistent file")
	}
}

func TestReaderFromFile_Valid(t *testing.T) {
	// This test expects a valid test file at testdata/test.jpg
	path := "../c2pa-rs/sdk/tests/fixtures/C.jpg"
	r, err := ReaderFromFile(path)
	assert.NotNil(t, r)
	assert.NotEmpty(t, r.Json())
	assert.Nil(t, err)
	if r != nil {
		r.Close()
	}
}
