package c2pa

import (
	"fmt"
	"runtime/cgo"
	"unsafe"
)

type Signer interface {
	Sign(input []byte, output []byte) (int, error)
	Alg() SigningAlg
	TimeStampUrl() string
	Certificates() string
}

// SignerInfo configures a built-in signer that signs locally with a private
// key. It maps to C2paSignerInfo and is consumed by ContextBuilder.SetSignerInfo.
type SignerInfo struct {
	// Alg is the signing algorithm.
	Alg SigningAlg
	// SignCert is the public certificate chain in PEM format.
	SignCert string
	// PrivateKey is the private key in PEM format.
	PrivateKey string
	// TimestampUrl is an optional RFC 3161 timestamp authority URL.
	TimestampUrl string
}

type NativeSigner struct {
	signer  Signer
	ptr     unsafe.Pointer
	handles []cgo.Handle
}

// goSignerCallback is invoked from the cgo signerCallback in native.go.
func goSignerCallback(handle uintptr, input, output []byte) (int, bool) {
	adapter, ok := cgo.Handle(handle).Value().(*NativeSigner)
	if !ok || adapter == nil || adapter.signer == nil {
		return 0, false
	}
	n, err := adapter.signer.Sign(input, output)
	if err != nil {
		return 0, false
	}
	return n, true
}

func (s NativeSigner) Close() {
	if s.ptr != nil {
		c2paFree(s.ptr)
		s.ptr = nil
	}
	for _, handle := range s.handles {
		if handle != 0 {
			handle.Delete()
		}
	}
	s.handles = nil
}

func (s *NativeSigner) ReserveSize() (int64, error) {
	size := c2paSignerReserveSize(s.ptr)
	if size == -1 {
		return 0, fmt.Errorf("failed to get signer reserve size: %s", c2paError())
	}
	return size, nil
}

// newSigner wraps a user-provided Signer in an adapter that exposes a
// C-callable signing callback. The returned adapter owns a C2paSigner and a
// cgo.Handle; both are released by Close.
func newSigner(signer Signer) (*NativeSigner, error) {
	s := &NativeSigner{
		signer: signer,
		ptr:    nil,
	}
	handle := cgo.NewHandle(s)
	s.handles = []cgo.Handle{handle}
	s.ptr = c2paSignerCreate(uintptr(handle), signer.Alg(), signer.TimeStampUrl(), signer.Certificates())
	if s.ptr == nil {
		handle.Delete()
		s.handles = nil
		return nil, fmt.Errorf("failed to create signer: %s", c2paError())
	}
	return s, nil
}

func NewSignerFromInfo(info SignerInfo) (*NativeSigner, error) {
	ptr := c2paSignerFromInfo(info.Alg, info.SignCert, info.PrivateKey, info.TimestampUrl)
	if ptr == nil {
		return nil, fmt.Errorf("failed to create signer from info: %s", c2paError())
	}
	return &NativeSigner{ptr: ptr}, nil
}

func NewIdentitySigner(c2paSigner Signer, identitySigner Signer, referencedAssertions []string, roles []string) (*NativeSigner, error) {
	is, err := newSigner(identitySigner)
	if err != nil {
		return nil, err
	}
	defer is.Close()
	cs, err := newSigner(c2paSigner)
	if err != nil {
		return nil, err
	}
	defer cs.Close()
	return newIdentitySigner(cs, is, referencedAssertions, roles)
}

// NewIdentitySigner combines two native signers into a single signer that
// emits both the C2PA claim signature and an X.509 identity assertion.
// On success, ownership of both input signers transfers to the returned signer.
func newIdentitySigner(c2paSigner *NativeSigner, identitySigner *NativeSigner, referencedAssertions []string, roles []string) (*NativeSigner, error) {
	if c2paSigner == nil || c2paSigner.ptr == nil {
		return nil, fmt.Errorf("c2pa signer is nil")
	}
	if identitySigner == nil || identitySigner.ptr == nil {
		return nil, fmt.Errorf("identity signer is nil")
	}

	ptr := c2paIdentitySignerCreate(c2paSigner.ptr, identitySigner.ptr, referencedAssertions, roles)
	if ptr == nil {
		return nil, fmt.Errorf("failed to create identity signer: %s", c2paError())
	}

	handles := append([]cgo.Handle{}, c2paSigner.handles...)
	handles = append(handles, identitySigner.handles...)

	c2paSigner.ptr = nil
	c2paSigner.handles = nil
	identitySigner.ptr = nil
	identitySigner.handles = nil

	return &NativeSigner{ptr: ptr, handles: handles}, nil
}

// ed25519Signer is a Signer that produces Ed25519 signatures via the native
// c2pa_ed25519_sign helper. Use NewEd25519Signer to construct one.
type ed25519Signer struct {
	privateKey   string
	certificates string
	timestampUrl string
}

// ed25519SignatureLen is the fixed size of an Ed25519 signature in bytes.
const ed25519SignatureLen = 64

// NewEd25519Signer returns a Signer that signs with the native c2pa Ed25519
// helper. privateKey is a PEM-encoded Ed25519 private key, certificates is the
// matching PEM certificate chain, and timestampUrl is an optional RFC 3161
// timestamp authority URL (pass "" to disable timestamping).
func NewEd25519Signer(privateKey, certificates, timestampUrl string) Signer {
	var s Signer = &ed25519Signer{
		privateKey:   privateKey,
		certificates: certificates,
		timestampUrl: timestampUrl,
	}
	return s
}

func (s *ed25519Signer) Sign(input []byte, output []byte) (int, error) {
	if len(output) < ed25519SignatureLen {
		return 0, fmt.Errorf("output buffer too small: need %d, got %d", ed25519SignatureLen, len(output))
	}
	sig := c2paEd25519Sign(input, s.privateKey)
	if sig == nil {
		return 0, fmt.Errorf("ed25519 sign failed: %s", c2paError())
	}
	copy(output[:ed25519SignatureLen], sig)
	return ed25519SignatureLen, nil
}

func (s *ed25519Signer) Alg() SigningAlg      { return SigningAlgEd25519 }
func (s *ed25519Signer) TimeStampUrl() string { return s.timestampUrl }
func (s *ed25519Signer) Certificates() string { return s.certificates }
