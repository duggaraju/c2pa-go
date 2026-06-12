//go:build linux && cgo
// +build linux,cgo

package c2pa

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime/cgo"
	"testing"

	"github.com/stretchr/testify/assert"
)

func fixturePath(parts ...string) string {
	all := append([]string{"..", "c2pa-rs"}, parts...)
	return filepath.Join(all...)
}

func TestVersionFormat(t *testing.T) {
	v := Version()
	assert.NotEmpty(t, v)
	assert.Regexp(t, regexp.MustCompile(`^c2pa-c-ffi/\d+\.\d+\.\d+\s+c2pa-rs/\d+\.\d+\.\d+$`), v)
}

func TestReaderAndBuilderNilContext(t *testing.T) {
	_, err := NewReader(nil)
	assert.Error(t, err)

	_, err = NewBuilder(nil)
	assert.Error(t, err)
}

func TestSettingsGuards(t *testing.T) {
	var settings Settings

	assert.Error(t, settings.UpdateFromString(`{}`, "json"))
	assert.Error(t, settings.SetValue("trust_list", `[]`))
	assert.Error(t, settings.UpdateFrom(nil))

	settings.Close()
	settings.Close()
}

func TestNewSettings(t *testing.T) {
	settings, err := NewSettings()
	assert.NoError(t, err)
	if settings != nil {
		defer settings.Close()
	}
}

func TestNewStream(t *testing.T) {
	file, err := os.CreateTemp(t.TempDir(), "stream-*.bin")
	assert.NoError(t, err)
	defer file.Close()

	stream, err := NewStream(file)
	assert.NoError(t, err)
	if stream != nil {
		stream.Close()
	}
}

func TestStreamCallbacks(t *testing.T) {
	readFile, err := os.CreateTemp(t.TempDir(), "read-*.bin")
	assert.NoError(t, err)
	_, err = readFile.WriteString("hello world")
	assert.NoError(t, err)
	_, err = readFile.Seek(0, io.SeekStart)
	assert.NoError(t, err)

	readHandle := cgo.NewHandle(Stream{file: readFile})
	defer readHandle.Delete()

	buf := make([]byte, 5)
	n := goStreamRead(uintptr(readHandle), buf)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte("hello"), buf)

	pos := goStreamSeek(uintptr(readHandle), 6, io.SeekStart)
	assert.Equal(t, int64(6), pos)

	buf = make([]byte, 5)
	n = goStreamRead(uintptr(readHandle), buf)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte("world"), buf)

	writeFile, err := os.CreateTemp(t.TempDir(), "write-*.bin")
	assert.NoError(t, err)
	writeHandle := cgo.NewHandle(Stream{file: writeFile})
	defer writeHandle.Delete()

	n = goStreamWrite(uintptr(writeHandle), []byte("abc"))
	assert.Equal(t, 3, n)
	assert.Equal(t, 0, goStreamFlush(uintptr(writeHandle)))
	assert.NoError(t, writeFile.Close())

	written, err := os.ReadFile(writeFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, []byte("abc"), written)
}

func TestBuilderPathWrappers(t *testing.T) {
	builder := &Builder{}

	missing := filepath.Join(t.TempDir(), "missing", "input.bin")
	assert.Error(t, builder.AddResourceFromFile("resource://missing", missing))
	assert.Error(t, builder.AddIngredientFromFile(`{"title":"missing"}`, missing))
	_, err := builder.FromArchiveFile(missing)
	assert.Error(t, err)

	output := filepath.Join(t.TempDir(), "missing", "archive.bin")
	assert.Error(t, builder.ToArchiveFile(output))
	assert.Error(t, builder.WriteIngredientArchiveFile("ingredient-id", output))
}

func TestReaderMetadataAndManifestParsing(t *testing.T) {
	ctx, err := NewContext()
	assert.NoError(t, err)
	defer ctx.Close()

	reader, err := NewReader(ctx)
	assert.NoError(t, err)
	defer reader.Close()

	asset := fixturePath("sdk", "tests", "fixtures", "C.jpg")
	assert.NoError(t, reader.WithFile(asset))
	assert.Empty(t, reader.RemoteUrl())
	assert.True(t, reader.IsEmbedded())
	assert.NotEmpty(t, reader.Json())
	assert.NotEmpty(t, reader.DetailedJson())

	manifest, err := reader.Manifest()
	assert.NoError(t, err)
	assert.NotNil(t, manifest)

	detailed, err := reader.DetailedManifest()
	assert.NoError(t, err)
	assert.NotNil(t, detailed)
}

type callbackSignerMock struct {
	alg      SigningAlg
	ts       string
	certs    string
	signFunc func(input, output []byte) (int, error)
}

func (s callbackSignerMock) Alg() SigningAlg      { return s.alg }
func (s callbackSignerMock) TimeStampUrl() string { return s.ts }
func (s callbackSignerMock) Certificates() string { return s.certs }
func (s callbackSignerMock) Sign(input []byte, output []byte) (int, error) {
	return s.signFunc(input, output)
}

type plainSignerMock struct{}

func (plainSignerMock) Alg() SigningAlg      { return SigningAlgEd25519 }
func (plainSignerMock) TimeStampUrl() string { return "" }
func (plainSignerMock) Certificates() string { return "" }

func TestGoSignerCallback(t *testing.T) {
	okSigner := callbackSignerMock{
		alg:   SigningAlgEd25519,
		ts:    "https://example.invalid/timestamp",
		certs: "cert-chain",
		signFunc: func(input, output []byte) (int, error) {
			copy(output, append([]byte("signed:"), input...))
			return len("signed:") + len(input), nil
		},
	}

	adapter := &NativeSigner{signer: okSigner}
	handle := cgo.NewHandle(adapter)
	defer handle.Delete()

	output := make([]byte, 32)
	n, ok := goSignerCallback(uintptr(handle), []byte("payload"), output)
	assert.True(t, ok)
	assert.Equal(t, len("signed:")+len("payload"), n)
	assert.True(t, bytes.HasPrefix(output, []byte("signed:payload")))

	badSigner := callbackSignerMock{
		alg:   SigningAlgEd25519,
		ts:    "https://example.invalid/timestamp",
		certs: "cert-chain",
		signFunc: func(input, output []byte) (int, error) {
			return 0, assert.AnError
		},
	}

	badAdapter := &NativeSigner{signer: badSigner}
	badHandle := cgo.NewHandle(badAdapter)
	defer badHandle.Delete()

	n, ok = goSignerCallback(uintptr(badHandle), []byte("payload"), output)
	assert.False(t, ok)
	assert.Zero(t, n)
}

func TestSignerAccessorsAndEd25519Signer(t *testing.T) {
	privateKeyPath := fixturePath("sdk", "tests", "fixtures", "crypto", "raw_signature", "ed25519.priv")
	certChainPath := fixturePath("sdk", "tests", "fixtures", "crypto", "raw_signature", "ed25519.pub")
	expectedSigPath := fixturePath("sdk", "tests", "fixtures", "crypto", "raw_signature", "ed25519.raw_sig")

	privateKey, err := os.ReadFile(privateKeyPath)
	assert.NoError(t, err)
	certChain, err := os.ReadFile(certChainPath)
	assert.NoError(t, err)
	expectedSig, err := os.ReadFile(expectedSigPath)
	assert.NoError(t, err)

	signer := NewEd25519Signer(string(privateKey), string(certChain), "https://timestamp.example.invalid")
	assert.Equal(t, SigningAlgEd25519, signer.Alg())
	assert.Equal(t, "https://timestamp.example.invalid", signer.TimeStampUrl())
	assert.Equal(t, string(certChain), signer.Certificates())

	_, err = signer.Sign([]byte("some sample content to sign"), make([]byte, ed25519SignatureLen-1))
	assert.Error(t, err)

	output := make([]byte, ed25519SignatureLen)
	n, err := signer.Sign([]byte("some sample content to sign"), output)
	assert.NoError(t, err)
	assert.Equal(t, ed25519SignatureLen, n)
	assert.True(t, bytes.Equal(expectedSig, output))
}

func TestSignerGuards(t *testing.T) {
	_, err := takeNativeSigner(nil)
	assert.Error(t, err)

	_, err = takeNativeSigner(plainSignerMock{})
	assert.Error(t, err)

	s := &signerFromInfo{info: SignerInfo{Alg: SigningAlgEd25519}}
	assert.Equal(t, SigningAlgEd25519, s.Alg())
	assert.Empty(t, s.TimeStampUrl())
	assert.Empty(t, s.Certificates())
	_, err = s.takeNativeSigner()
	assert.Error(t, err)
	_, err = s.ReserveSize()
	assert.Error(t, err)
}

func TestBuilderSignWrapper(t *testing.T) {
	manifestJSON, err := os.ReadFile(fixturePath("sdk", "tests", "fixtures", "simple_manifest.json"))
	assert.NoError(t, err)

	signCert, err := os.ReadFile(fixturePath("sdk", "tests", "fixtures", "certs", "ps256.pub"))
	assert.NoError(t, err)

	privateKey, err := os.ReadFile(fixturePath("sdk", "tests", "fixtures", "certs", "ps256.pem"))
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

	builder, err := NewBuilder(ctx)
	assert.NoError(t, err)
	defer builder.Close()

	builder, err = builder.WithDefinition(string(manifestJSON))
	assert.NoError(t, err)

	input := fixturePath("sdk", "tests", "fixtures", "C.jpg")
	output := filepath.Join(t.TempDir(), "signed.jpg")

	manifest, err := builder.Sign(input, output, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, manifest)

	info, err := os.Stat(output)
	assert.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))
}
