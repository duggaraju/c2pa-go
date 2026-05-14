package c2pa

// #include <c2pa.h>
import "C"

import (
	"encoding/json"
	"fmt"
	"unsafe"

	schema "github.com/duggaraju/c2pa-go/c2pa/schema"
)

// Settings wraps a C2paSettings*. It is used to configure a ContextBuilder via
// ContextBuilder.SetSettings.
type Settings struct {
	ptr *C.C2paSettings
}

// NewSettings creates a new Settings with defaults.
func NewSettings() (*Settings, error) {
	ptr := C.c2pa_settings_new()
	if ptr == nil {
		return nil, fmt.Errorf("failed to create c2pa settings: %s", c2paError())
	}
	return &Settings{ptr: ptr}, nil
}

// Close releases the underlying C settings. Safe to call multiple times.
func (s *Settings) Close() {
	if s.ptr != nil {
		C.c2pa_free(unsafe.Pointer(s.ptr))
		s.ptr = nil
	}
}

// UpdateFromString loads settings from a JSON or TOML string. The format
// argument must be "json" or "toml".
func (s *Settings) UpdateFromString(content, format string) error {
	if s.ptr == nil {
		return fmt.Errorf("settings is closed")
	}
	cContent := C.CString(content)
	defer C.free(unsafe.Pointer(cContent))
	cFormat := C.CString(format)
	defer C.free(unsafe.Pointer(cFormat))
	if rc := C.c2pa_settings_update_from_string(s.ptr, cContent, cFormat); rc != 0 {
		return fmt.Errorf("failed to update settings: %s", c2paError())
	}
	return nil
}

// SetValue sets a single configuration value using dot notation. The value
// must be a JSON-encoded scalar or array (e.g. "true", "42", "\"ps256\"").
func (s *Settings) SetValue(path, jsonValue string) error {
	if s.ptr == nil {
		return fmt.Errorf("settings is closed")
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	cValue := C.CString(jsonValue)
	defer C.free(unsafe.Pointer(cValue))
	if rc := C.c2pa_settings_set_value(s.ptr, cPath, cValue); rc != 0 {
		return fmt.Errorf("failed to set settings value: %s", c2paError())
	}
	return nil
}

// UpdateFrom applies a typed schema.Settings value by marshaling it to
// JSON and calling UpdateFromString.
func (s *Settings) UpdateFrom(settings *schema.Settings) error {
	if settings == nil {
		return fmt.Errorf("settings is nil")
	}
	data, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}
	return s.UpdateFromString(string(data), "json")
}
