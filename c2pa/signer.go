package c2pa

import (
	"fmt"
	"runtime/cgo"
	"unsafe"
)

type Signer interface {
	Alg() SigningAlg
	TimeStampUrl() string
	Certificates() string
}

// CallbackSigner extends Signer with a callback-based Sign method.
// Use this for Go-implemented signers that receive input bytes to sign.
type CallbackSigner interface {
	Signer
	Sign(input []byte, output []byte) (int, error)
}

// nativeBackedSigner is implemented by Signer implementations that already
// own a native C2PA signer and can transfer ownership to the caller.
type nativeBackedSigner interface {
	Signer
	takeNativeSigner() (*NativeSigner, error)
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
	signer  CallbackSigner
	ptr     unsafe.Pointer
	handles []cgo.Handle
}

// signerFromInfo is a Signer backed by a native signer created from SignerInfo.
type signerFromInfo struct {
	info   SignerInfo
	native *NativeSigner
}

// identitySigner is a Signer backed by a native identity signer.
type identitySigner struct {
	claim    Signer
	identity Signer
	native   *NativeSigner
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

// newSigner wraps a user-provided CallbackSigner in an adapter that exposes a
// C-callable signing callback. The returned adapter owns a C2paSigner and a
// cgo.Handle; both are released by Close.
func newSigner(signer CallbackSigner) (*NativeSigner, error) {
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

func nativeSignerFromInfo(info SignerInfo) (*NativeSigner, error) {
	ptr := c2paSignerFromInfo(info.Alg, info.SignCert, info.PrivateKey, info.TimestampUrl)
	if ptr == nil {
		return nil, fmt.Errorf("failed to create signer from info: %s", c2paError())
	}
	return &NativeSigner{ptr: ptr}, nil
}

// NewSignerFromInfo creates a native-backed Signer from certificate/private-key
// material. The returned signer can be passed directly to SetSigner or
// NewIdentitySigner.
func NewSignerFromInfo(info SignerInfo) (Signer, error) {
	native, err := nativeSignerFromInfo(info)
	if err != nil {
		return nil, err
	}
	return &signerFromInfo{info: info, native: native}, nil
}

func takeNativeSigner(signer Signer) (*NativeSigner, error) {
	if signer == nil {
		return nil, fmt.Errorf("signer is nil")
	}
	if nativeBacked, ok := signer.(nativeBackedSigner); ok {
		return nativeBacked.takeNativeSigner()
	}
	if callbackSigner, ok := signer.(CallbackSigner); ok {
		return newSigner(callbackSigner)
	}
	return nil, fmt.Errorf("signer does not implement CallbackSigner and is not native-backed")
}

// NewIdentitySigner combines two Signers into a single Signer that emits both
// the C2PA claim signature and an X.509 identity assertion.
func NewIdentitySigner(c2paSigner Signer, idSigner Signer, referencedAssertions []string, roles []string) (Signer, error) {
	is, err := takeNativeSigner(idSigner)
	if err != nil {
		return nil, err
	}
	cs, err := takeNativeSigner(c2paSigner)
	if err != nil {
		is.Close()
		return nil, err
	}

	native, err := newIdentitySigner(cs, is, referencedAssertions, roles)
	if err != nil {
		cs.Close()
		is.Close()
		return nil, err
	}
	return &identitySigner{claim: c2paSigner, identity: idSigner, native: native}, nil
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

func (s *signerFromInfo) Alg() SigningAlg {
	if s == nil {
		return ""
	}
	return s.info.Alg
}

func (s *signerFromInfo) TimeStampUrl() string {
	if s == nil {
		return ""
	}
	return s.info.TimestampUrl
}

func (s *signerFromInfo) Certificates() string {
	if s == nil {
		return ""
	}
	return s.info.SignCert
}

func (s *signerFromInfo) takeNativeSigner() (*NativeSigner, error) {
	if s == nil || s.native == nil || s.native.ptr == nil {
		return nil, fmt.Errorf("signer is nil")
	}
	native := s.native
	s.native = nil
	return native, nil
}

func (s *signerFromInfo) ReserveSize() (int64, error) {
	if s == nil || s.native == nil {
		return 0, fmt.Errorf("signer is nil")
	}
	return s.native.ReserveSize()
}

func (s *signerFromInfo) Close() {
	if s == nil || s.native == nil {
		return
	}
	s.native.Close()
	s.native = nil
}

func (s *identitySigner) Alg() SigningAlg {
	if s == nil || s.claim == nil {
		return ""
	}
	return s.claim.Alg()
}

func (s *identitySigner) TimeStampUrl() string {
	if s == nil || s.claim == nil {
		return ""
	}
	return s.claim.TimeStampUrl()
}

func (s *identitySigner) Certificates() string {
	if s == nil || s.claim == nil {
		return ""
	}
	return s.claim.Certificates()
}

func (s *identitySigner) takeNativeSigner() (*NativeSigner, error) {
	if s == nil || s.native == nil || s.native.ptr == nil {
		return nil, fmt.Errorf("identity signer is nil")
	}
	native := s.native
	s.native = nil
	return native, nil
}

func (s *identitySigner) ReserveSize() (int64, error) {
	if s == nil || s.native == nil {
		return 0, fmt.Errorf("identity signer is nil")
	}
	return s.native.ReserveSize()
}

func (s *identitySigner) Close() {
	if s == nil || s.native == nil {
		return
	}
	s.native.Close()
	s.native = nil
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
func NewEd25519Signer(privateKey, certificates, timestampUrl string) CallbackSigner {
	var s CallbackSigner = &ed25519Signer{
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
