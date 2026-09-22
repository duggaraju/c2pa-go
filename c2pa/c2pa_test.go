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
	"github.com/stretchr/testify/require"
)

func TestC2paVersion(t *testing.T) {
	v := Version()
	assert.NotEmpty(t, v)
	assert.Regexp(t, regexp.MustCompile(`^c2pa-c-ffi/\d+\.\d+\.\d+\s+c2pa-rs/\d+\.\d+\.\d+$`), v)
}

func TestContextBuilderSetProgressCallback_Invoked(t *testing.T) {
	signCert, err := os.ReadFile("../c2pa-rs/sdk/tests/fixtures/certs/ps256.pub")
	require.NoError(t, err)

	privateKey, err := os.ReadFile("../c2pa-rs/sdk/tests/fixtures/certs/ps256.pem")
	require.NoError(t, err)

	ctxBuilder, err := NewContextBuilder()
	require.NoError(t, err)
	defer ctxBuilder.Close()

	var callbacks atomic.Int32
	err = ctxBuilder.SetProgressCallback(ProgressFunc(func(phase ProgressPhase, step uint32, total uint32) bool {
		callbacks.Add(1)
		return true
	}))
	require.NoError(t, err)

	err = ctxBuilder.SetSignerInfo(SignerInfo{
		Alg:        "ps256",
		SignCert:   string(signCert),
		PrivateKey: string(privateKey),
	})
	require.NoError(t, err)

	ctx, err := ctxBuilder.Build()
	require.NoError(t, err)
	defer ctx.Close()

	b, err := NewBuilder(ctx)
	require.NoError(t, err)
	defer b.Close()

	b, err = b.WithDefinition(signingManifestJSON)
	require.NoError(t, err)
	require.NoError(t, b.SetIntent(IntentEdit, SourceEmpty))

	input := "../c2pa-rs/sdk/tests/fixtures/C.jpg"
	output := filepath.Join(t.TempDir(), "signed.jpg")

	_, err = b.SignWithContext(input, output)
	require.NoError(t, err)
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
	signCert, err := os.ReadFile("../c2pa-rs/sdk/tests/fixtures/certs/ps256.pub")
	require.NoError(t, err)

	privateKey, err := os.ReadFile("../c2pa-rs/sdk/tests/fixtures/certs/ps256.pem")
	require.NoError(t, err)

	ctxBuilder, err := NewContextBuilder()
	require.NoError(t, err)
	defer ctxBuilder.Close()

	err = ctxBuilder.SetSignerInfo(SignerInfo{
		Alg:        "ps256",
		SignCert:   string(signCert),
		PrivateKey: string(privateKey),
	})
	require.NoError(t, err)

	ctx, err := ctxBuilder.Build()
	require.NoError(t, err)
	defer ctx.Close()

	b, err := NewBuilder(ctx)
	require.NoError(t, err)
	defer b.Close()

	b, err = b.WithDefinition(signingManifestJSON)
	require.NoError(t, err)
	require.NoError(t, b.SetIntent(IntentEdit, SourceEmpty))

	input := "../c2pa-rs/sdk/tests/fixtures/C.jpg"
	output := filepath.Join(t.TempDir(), "signed.jpg")

	manifest, err := b.SignWithContext(input, output)
	require.NoError(t, err)
	assert.NotEmpty(t, manifest)

	info, err := os.Stat(output)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))

	r, err := NewReader(ctx)
	require.NoError(t, err)
	defer r.Close()
	f, err := os.Open(output)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, f.Close())
	}()
	require.NoError(t, r.WithStream("jpg", f))

	store, err := r.Manifest()
	require.NoError(t, err)
	require.NotNil(t, store.ActiveManifest)
	active, ok := store.Manifests[*store.ActiveManifest]
	require.True(t, ok)
	require.Len(t, active.Ingredients, 1)
	require.NotNil(t, active.Ingredients[0].Relationship)
	assert.Equal(t, "parentOf", string(*active.Ingredients[0].Relationship))
	assert.NotEmpty(t, active.Assertions)
}
