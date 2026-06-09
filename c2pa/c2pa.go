package c2pa

// Version returns the version string from the c2pa library.
func Version() string {
	return c2paVersion()
}

// SetError sets the c2pa thread-local last-error string. Useful when
// implementing callbacks (signers, http resolvers, progress) so that the
// SDK can surface the failure reason via c2pa_error.
func SetError(error string) {
	c2paErrorSetLast(error)
}

// BuilderSupportedMimeTypes returns the MIME types supported by the C2PA builder.
func BuilderSupportedMimeTypes() []string {
	return c2paBuilderSupportedMimeTypes()
}

// ReaderSupportedMimeTypes returns the MIME types supported by the C2PA reader.
func ReaderSupportedMimeTypes() []string {
	return c2paReaderSupportedMimeTypes()
}
