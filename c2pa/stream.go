package c2pa

import (
	"fmt"
	"io"
	"os"
	"runtime/cgo"
	"unsafe"
)

type Stream struct {
	file   *os.File
	ptr    unsafe.Pointer
	handle cgo.Handle
}

// goStreamRead is invoked from the cgo streamRead callback in native.go.
// In a build without cgo it is unreachable but must remain compilable.
func goStreamRead(handle uintptr, buf []byte) int {
	stream := cgo.Handle(handle).Value().(Stream)
	n, err := stream.file.Read(buf)
	if err != nil {
		if err == io.EOF {
			return 0
		}
		return -1
	}
	return n
}

func goStreamSeek(handle uintptr, offset int64, mode int) int64 {
	stream := cgo.Handle(handle).Value().(Stream)
	n, err := stream.file.Seek(offset, mode)
	if err != nil {
		return -1
	}
	return n
}

func goStreamWrite(handle uintptr, buf []byte) int {
	stream := cgo.Handle(handle).Value().(Stream)
	n, _ := stream.file.Write(buf)
	return n
}

func goStreamFlush(handle uintptr) int {
	stream := cgo.Handle(handle).Value().(Stream)
	if err := stream.file.Sync(); err != nil {
		return -1
	}
	return 0
}

// NewStream creates a new Stream.
func NewStream(file *os.File) (*Stream, error) {
	stream := Stream{
		ptr:  nil,
		file: file,
	}

	stream.handle = cgo.NewHandle(stream)
	stream.ptr = c2paCreateStream(uintptr(stream.handle))
	if stream.ptr == nil {
		err := c2paError()
		stream.handle.Delete()
		return nil, fmt.Errorf("failed to create stream: %s", err)
	}

	return &stream, nil
}

func (s *Stream) Close() {
	if s.ptr != nil {
		c2paReleaseStream(s.ptr)
		s.ptr = nil
	}
	if s.handle != 0 {
		s.handle.Delete()
		s.handle = 0
	}
}
