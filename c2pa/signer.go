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
