package lib

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

type SignerAdapter struct {
	signer Signer
	ptr    *C.C2paSigner
	handle cgo.Handle
}

//export GoSignerCallback
func GoSignerCallback(context C.uintptr_t, input *C.uint8_t, input_size C.uintptr_t, output *C.uint8_t, output_size C.uintptr_t) C.intptr_t {
	handle := cgo.Handle(context)
	signerAdapter := handle.Value().(SignerAdapter)

	in := C.GoBytes(unsafe.Pointer(input), C.int(input_size))
	out := C.GoBytes(unsafe.Pointer(output), C.int(output_size))

	n, err := signerAdapter.signer.Sign(in, out)
	if err != nil {
		return C.intptr_t(-1)
	}
	return C.intptr_t(n)
}

func (s *SignerAdapter) Close() {
	C.c2pa_signer_free(s.ptr)
	s.handle.Delete()
}

func (s *SignerAdapter) Sign(input []byte, output []byte) (int, error) {
	if s.signer != nil {
		return s.signer.Sign(input, output)
	}
	return -1, nil
}

func NewSigner(signer Signer) (*SignerAdapter, error) {
	s := &SignerAdapter{
		signer: signer,
		ptr:    nil,
	}
	s.handle = cgo.NewHandle(s)
	taUrl := C.CString(signer.TimeStampUrl())
	defer C.free(unsafe.Pointer(taUrl))

	certificates := C.CString(signer.Certificates())
	defer C.free(unsafe.Pointer(certificates))

	s.ptr = C.create_signer(C.uintptr_t(s.handle), C.C2paSigningAlg(signer.Alg()), taUrl, certificates)
	if s.ptr == nil {
		return nil, fmt.Errorf("failed to create signer: %s", C2paError())
	}
	return s, nil
}
