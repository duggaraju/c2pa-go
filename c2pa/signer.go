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

type signerAdapter struct {
	signer Signer
	ptr    *C.C2paSigner
	handle cgo.Handle
}

//export signerCallback
func signerCallback(context C.uintptr_t, input *C.uint8_t, input_size C.uintptr_t, output *C.uint8_t, output_size C.uintptr_t) C.intptr_t {
	handle := cgo.Handle(context)
	adapter, ok := handle.Value().(*signerAdapter)
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

func (s *signerAdapter) Close() {
	if s.ptr != nil {
		C.c2pa_free(unsafe.Pointer(s.ptr))
		s.ptr = nil
	}
	s.handle.Delete()
}

func (s *signerAdapter) Sign(input []byte, output []byte) (int, error) {
	if s.signer != nil {
		return s.signer.Sign(input, output)
	}
	return -1, nil
}

// newSigner wraps a user-provided Signer in an adapter that exposes a
// C-callable signing callback. The returned adapter owns a C2paSigner and a
// cgo.Handle; both are released by Close.
func newSigner(signer Signer) (*signerAdapter, error) {
	s := &signerAdapter{
		signer: signer,
		ptr:    nil,
	}
	s.handle = cgo.NewHandle(s)
	taUrl := C.CString(signer.TimeStampUrl())
	defer C.free(unsafe.Pointer(taUrl))

	certs := signer.Certificates()
	certificates := C.CString(certs)
	defer C.free(unsafe.Pointer(certificates))

	s.ptr = C.create_signer(C.uintptr_t(s.handle), C.C2paSigningAlg(signer.Alg()), taUrl, certificates)
	if s.ptr == nil {
		return nil, fmt.Errorf("failed to create signer: %s", c2paError())
	}
	return s, nil
}
