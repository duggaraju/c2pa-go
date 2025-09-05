package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCpaVersion(t *testing.T) {
	v := CpaVersion()
	expected := "c2pa-c-ffi/0.61.0 c2pa-rs/0.61.0"
	assert.Equal(t, expected, v)
}

func TestReaderFromFile_NotFound(t *testing.T) {
	_, err := ReaderFromFile("/nonexistent/file/path.jpg")
	if err == nil {
		t.Error("ReaderFromFile should fail for nonexistent file")
	}
}

func TestReaderFromFile_Valid(t *testing.T) {
	// This test expects a valid test file at testdata/test.jpg
	path := "../c2pa-rs/sdk/tests/fixtures/C.jpg"
	r, err := ReaderFromFile(path)
	assert.NotNil(t, r)
	assert.NotEmpty(t, r.Json())
	assert.Nil(t, err)
	if r != nil {
		r.Close()
	}
}
