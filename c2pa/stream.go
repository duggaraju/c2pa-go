package c2pa

// #include "c2pa_helper.h"
import "C"

import (
	"fmt"
	"io"
	"os"
	"runtime/cgo"
	"unsafe"
)

type Stream struct {
	file   *os.File
	ptr    *C.C2paStream
	handle cgo.Handle
}

//export streamRead
func streamRead(context C.uintptr_t, buffer *C.uint8_t, size C.intptr_t) C.intptr_t {
	handle := cgo.Handle(context)
	stream := handle.Value().(Stream)
	slice := unsafe.Slice((*byte)(buffer), (int)(size))
	n, err := stream.file.Read(slice)
	if err != nil {
		if err == io.EOF {
			return C.intptr_t(0) // EOF is not an error for Read
		}
		return C.intptr_t(-1)
	}
	return C.intptr_t(n)
}

//export streamSeek
func streamSeek(context C.uintptr_t, offset C.intptr_t, mode C.C2paSeekMode) C.intptr_t {
	handle := cgo.Handle(context)
	stream := handle.Value().(Stream)
	n, err := stream.file.Seek(int64(offset), int(mode))
	if err != nil {
		return C.intptr_t(-1)
	}
	return C.intptr_t(n)
}

//export streamWrite
func streamWrite(context C.uintptr_t, buffer *C.uint8_t, size C.intptr_t) C.intptr_t {
	handle := cgo.Handle(context)
	stream := handle.Value().(Stream)
	slice := unsafe.Slice((*byte)(buffer), (int)(size))
	n, err := stream.file.Write(slice)
	if err != nil {
		return C.intptr_t(n)
	}
	return C.intptr_t(n)
}

//export streamFlush
func streamFlush(context C.uintptr_t) C.intptr_t {
	handle := cgo.Handle(context)
	stream := handle.Value().(Stream)
	err := stream.file.Sync()
	if err == nil {
		return C.intptr_t(0)
	} else {
		return C.intptr_t(-1)
	}
}

// NewStream creates a new Stream.
func NewStream(file *os.File) (*Stream, error) {

	stream := Stream{
		ptr:  nil,
		file: file,
	}

	stream.handle = cgo.NewHandle(stream)
	stream.ptr = C.create_stream(C.uintptr_t(stream.handle))
	if stream.ptr == nil {
		err := c2paError()
		return nil, fmt.Errorf("failed to create stream: %s", err)
	}

	return &stream, nil
}

func (s *Stream) Close() {
	C.c2pa_release_stream(s.ptr)
	s.handle.Delete()
}
