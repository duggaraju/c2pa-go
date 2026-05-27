//go:build linux && cgo
// +build linux,cgo

package c2pa

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestC2paVersion(t *testing.T) {
	v := C2paVersion()
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

	ctx, ctxErr := NewContext()
	assert.NoError(t, ctxErr)
	defer ctx.Close()

	_, readerErr := ReaderFromFile(ctx, path)
	assert.Error(t, readerErr)

	lastErr := c2paError()
	assert.NotEmpty(t, lastErr)
}

func TestReaderFromFile_NotFound(t *testing.T) {
	ctx, err := NewContext()
	assert.NoError(t, err)
	defer ctx.Close()

	_, err = ReaderFromFile(ctx, "/nonexistent/file/path.jpg")
	if err == nil {
		t.Error("ReaderFromFile should fail for nonexistent file")
	}
}

func TestReaderFromFile_Valid(t *testing.T) {
	ctx, ctxErr := NewContext()
	assert.NoError(t, ctxErr)
	defer ctx.Close()

	// This test expects a valid test file at testdata/test.jpg
	path := "../c2pa-rs/sdk/tests/fixtures/C.jpg"
	r, err := ReaderFromFile(ctx, path)
	assert.NotNil(t, r)
	assert.NotEmpty(t, r.Json())
	assert.Nil(t, err)
	if r != nil {
		r.Close()
	}
}

func TestBuilderSignWithContext_Valid(t *testing.T) {
	manifestJSON, err := os.ReadFile("../c2pa-rs/sdk/tests/fixtures/simple_manifest.json")
	assert.NoError(t, err)

	signCert, err := os.ReadFile("../c2pa-rs/sdk/tests/fixtures/certs/ps256.pub")
	assert.NoError(t, err)

	privateKey, err := os.ReadFile("../c2pa-rs/sdk/tests/fixtures/certs/ps256.pem")
	assert.NoError(t, err)

	ctxBuilder, err := NewContextBuilder()
	assert.NoError(t, err)
	defer ctxBuilder.Close()

	err = ctxBuilder.SetSignerInfo(SignerInfo{
		Alg:        "ps256",
		SignCert:   string(signCert),
		PrivateKey: string(privateKey),
	})
	assert.NoError(t, err)

	ctx, err := ctxBuilder.Build()
	assert.NoError(t, err)
	defer ctx.Close()

	b, err := NewBuilder(ctx)
	assert.NoError(t, err)
	defer b.Close()

	b, err = b.WithDefinition(string(manifestJSON))
	assert.NoError(t, err)

	input := "../c2pa-rs/sdk/tests/fixtures/C.jpg"
	output := filepath.Join(t.TempDir(), "signed.jpg")

	manifest, err := b.SignWithContext(input, output)
	assert.NoError(t, err)
	assert.NotEmpty(t, manifest)

	info, err := os.Stat(output)
	assert.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))

	r, err := ReaderFromFile(ctx, output)
	assert.NoError(t, err)
	if r != nil {
		defer r.Close()
		assert.NotEmpty(t, r.Json())
	}
}
