package c2pa

// #include "c2pa_helper.h"
import "C"
import (
	"fmt"
	"runtime/cgo"
	"unsafe"
)

type SigningAlg C.C2paSigningAlg

const (
	SigningAlgPs256       SigningAlg = C.Ps256
	SigningAlgPs384       SigningAlg = C.Ps384
	SigningAlgPs512       SigningAlg = C.Ps512
	SigningAlgEs256       SigningAlg = C.Es256
	SigningAlgEs384       SigningAlg = C.Es384
	SigningAlgEs512       SigningAlg = C.Es512
	C2paSigningAlgEd25519 SigningAlg = C.Ed25519
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
	// Alg is the signing algorithm name (e.g. "ps256", "es256").
	Alg string
	// SignCert is the public certificate chain in PEM format.
	SignCert string
	// PrivateKey is the private key in PEM format.
	PrivateKey string
	// TimestampUrl is an optional RFC 3161 timestamp authority URL.
	TimestampUrl string
}

type NativeSigner struct {
	signer Signer
	ptr    *C.C2paSigner
	handles []cgo.Handle
}

//export signerCallback
func signerCallback(context C.uintptr_t, input *C.uint8_t, input_size C.uintptr_t, output *C.uint8_t, output_size C.uintptr_t) C.intptr_t {
	handle := cgo.Handle(context)
	adapter, ok := handle.Value().(*NativeSigner)
	if !ok || adapter == nil || adapter.signer == nil {
		return C.intptr_t(-1)
	}

	in := unsafe.Slice((*byte)(unsafe.Pointer(input)), int(input_size))
	out := unsafe.Slice((*byte)(unsafe.Pointer(output)), int(output_size))

	n, err := adapter.signer.Sign(in, out)
	if err != nil {
		return C.intptr_t(-1)
	}
	return C.intptr_t(n)
}

func (s NativeSigner) Close() {
	if s.ptr != nil {
		C.c2pa_free(unsafe.Pointer(s.ptr))
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
	size := int64(C.c2pa_signer_reserve_size(s.ptr))
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
	taUrl := C.CString(signer.TimeStampUrl())
	defer C.free(unsafe.Pointer(taUrl))

	certs := signer.Certificates()
	certificates := C.CString(certs)
	defer C.free(unsafe.Pointer(certificates))

	s.ptr = C.create_signer(C.uintptr_t(handle), C.C2paSigningAlg(signer.Alg()), taUrl, certificates)
	if s.ptr == nil {
		return nil, fmt.Errorf("failed to create signer: %s", c2paError())
	}
	return s, nil
}

func NewSignerFromInfo(info SignerInfo) (*NativeSigner, error) {
	cAlg := C.CString(info.Alg)
	defer C.free(unsafe.Pointer(cAlg))
	cCert := C.CString(info.SignCert)
	defer C.free(unsafe.Pointer(cCert))
	cKey := C.CString(info.PrivateKey)
	defer C.free(unsafe.Pointer(cKey))
	var cTaUrl *C.char
	if info.TimestampUrl != "" {
		cTaUrl = C.CString(info.TimestampUrl)
		defer C.free(unsafe.Pointer(cTaUrl))
	}

	cInfo := C.C2paSignerInfo {
		cAlg,
		cCert,
		cKey,
		cTaUrl,
	}
	ptr := C.c2pa_signer_from_info(&cInfo)
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

	crefs, freeRefs := cStringArray(referencedAssertions)
	defer freeRefs()
	rolesPtr, freeRoles := cStringArray(roles)
	defer freeRoles()

	ptr := C.c2pa_identity_signer_create(c2paSigner.ptr, identitySigner.ptr, crefs, rolesPtr)
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
	cKey := C.CString(s.privateKey)
	defer C.free(unsafe.Pointer(cKey))

	var inPtr *C.uchar
	if len(input) > 0 {
		inPtr = (*C.uchar)(unsafe.Pointer(&input[0]))
	}

	sig := C.c2pa_ed25519_sign(inPtr, C.uintptr_t(len(input)), cKey)
	if sig == nil {
		return 0, fmt.Errorf("ed25519 sign failed: %s", c2paError())
	}
	defer C.c2pa_free(unsafe.Pointer(sig))

	src := unsafe.Slice((*byte)(unsafe.Pointer(sig)), ed25519SignatureLen)
	copy(output[:ed25519SignatureLen], src)
	return ed25519SignatureLen, nil
}

func (s *ed25519Signer) Alg() SigningAlg      { return C2paSigningAlgEd25519 }
func (s *ed25519Signer) TimeStampUrl() string { return s.timestampUrl }
func (s *ed25519Signer) Certificates() string { return s.certificates }

func cStringArray(values []string) (**C.char, func()) {
	if len(values) == 0 {
		return nil, func() {}
	}
	ptrs := make([]*C.char, len(values)+1)
	for i, value := range values {
		ptrs[i] = C.CString(value)
	}
	return &ptrs[0], func() {
		for _, ptr := range ptrs[:len(values)] {
			C.free(unsafe.Pointer(ptr))
		}
	}
}
