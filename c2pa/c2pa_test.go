//go:build linux && cgo
// +build linux,cgo

package c2pa

import (
	"os"
	"path/filepath"
	"regexp"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestC2paVersion(t *testing.T) {
	v := Version()
	assert.NotEmpty(t, v)
	assert.Regexp(t, regexp.MustCompile(`^c2pa-c-ffi/\d+\.\d+\.\d+\s+c2pa-rs/\d+\.\d+\.\d+$`), v)
}

func TestContextBuilderSetProgressCallback_Invoked(t *testing.T) {
	manifestJSON, err := os.ReadFile("../c2pa-rs/sdk/tests/fixtures/simple_manifest.json")
	assert.NoError(t, err)

	signCert, err := os.ReadFile("../c2pa-rs/sdk/tests/fixtures/certs/ps256.pub")
	assert.NoError(t, err)

	privateKey, err := os.ReadFile("../c2pa-rs/sdk/tests/fixtures/certs/ps256.pem")
	assert.NoError(t, err)

	ctxBuilder, err := NewContextBuilder()
	assert.NoError(t, err)
	defer ctxBuilder.Close()

	var callbacks atomic.Int32
	err = ctxBuilder.SetProgressCallback(ProgressFunc(func(phase ProgressPhase, step uint32, total uint32) bool {
		callbacks.Add(1)
		return true
	}))
	assert.NoError(t, err)

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

	_, err = b.SignWithContext(input, output)
	assert.NoError(t, err)
	assert.Greater(t, callbacks.Load(), int32(0))
}

func TestSignerInfoReserveSize(t *testing.T) {
	signCert, err := os.ReadFile("../c2pa-rs/sdk/tests/fixtures/certs/ps256.pub")
	assert.NoError(t, err)

	privateKey, err := os.ReadFile("../c2pa-rs/sdk/tests/fixtures/certs/ps256.pem")
	assert.NoError(t, err)

	signer, err := NewSignerFromInfo(SignerInfo{
		Alg:        "ps256",
		SignCert:   string(signCert),
		PrivateKey: string(privateKey),
	})
	assert.NoError(t, err)
	reserveSizer, ok := signer.(interface{ ReserveSize() (int64, error) })
	assert.True(t, ok)
	size, err := reserveSizer.ReserveSize()
	assert.NoError(t, err)
	assert.Greater(t, size, int64(0))
}

func TestNewIdentitySigner(t *testing.T) {
	signCert, err := os.ReadFile("../c2pa-rs/sdk/tests/fixtures/certs/ps256.pub")
	assert.NoError(t, err)

	privateKey, err := os.ReadFile("../c2pa-rs/sdk/tests/fixtures/certs/ps256.pem")
	assert.NoError(t, err)

	claimSigner, err := NewSignerFromInfo(SignerInfo{
		Alg:        "ps256",
		SignCert:   string(signCert),
		PrivateKey: string(privateKey),
	})
	assert.NoError(t, err)

	identitySigner, err := NewSignerFromInfo(SignerInfo{
		Alg:        "ps256",
		SignCert:   string(signCert),
		PrivateKey: string(privateKey),
	})
	assert.NoError(t, err)

	signer, err := NewIdentitySigner(claimSigner, identitySigner, []string{"c2pa.actions"}, []string{"author"})
	assert.NoError(t, err)
	if signer != nil {
		closer, ok := signer.(interface{ Close() })
		if ok {
			defer closer.Close()
		}
		reserveSizer, ok := signer.(interface{ ReserveSize() (int64, error) })
		assert.True(t, ok)
		size, reserveErr := reserveSizer.ReserveSize()
		assert.NoError(t, reserveErr)
		assert.Greater(t, size, int64(0))
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

	r, err := NewReader(ctx)
	assert.NoError(t, err)
	if r != nil {
		defer r.Close()
		f, openErr := os.Open(output)
		assert.NoError(t, openErr)
		if f != nil {
			defer func() {
				assert.NoError(t, f.Close())
			}()
			assert.NoError(t, r.WithStream("jpg", f))
		}
		assert.NotEmpty(t, r.Json())
	}
}
