package converter

import "fmt"

// SupportedMimeTypes is the list of all mimetypes that can be converted.
var SupportedMimeTypes []string

// OutputMimeTypes is the set of output formats.
var OutputMimeTypes = make(map[string]struct{})

// MediaConverter is the exported converter that can convert any mimetype in `SupportedMimeTypes`.
type MediaConverter struct {
	converters map[string]converter
}

// converter describes the behavior of all converters (HEIC, quicktime, etc).
type converter interface {
	// Convert copies `src` to `dst`, converting between the two files if necessary.
	// `dst` may not yet exist after this function exits.
	// Returns the full path of the destination after adding any suffixes.
	// `dst` should not have any file extensions, as they will be added by the converter.
	Convert(src, dst string) (string, error)

	// Complete any remaining conversion tasks, blocking until all are done.
	Finish() error
}

var registeredConverters = make(map[string]func() converter)

// registerConverter registers a converter for a specific mimetype.
func registerConverter(inmime, outmime string, c func() converter) int {
	if _, alreadyRegistered := registeredConverters[inmime]; alreadyRegistered {
		panic(fmt.Errorf("converter already exists for mimetype %s", inmime))
	}
	registeredConverters[inmime] = c

	SupportedMimeTypes = append(SupportedMimeTypes, inmime)
	OutputMimeTypes[outmime] = struct{}{}

	return 0
}

func NewMediaConverter() *MediaConverter {
	m := &MediaConverter{converters: make(map[string]converter)}
	for mimetype, c := range registeredConverters {
		m.converters[mimetype] = c()
	}

	return m
}

func (m *MediaConverter) Convert(src, dst, srcMimeType string) (string, error) {
	c, converterExists := m.converters[srcMimeType]
	if !converterExists {
		return "", fmt.Errorf("mimetype not supported: %s", srcMimeType)
	}

	return c.Convert(src, dst)
}

func (m *MediaConverter) Finish() error {
	for mimetype, c := range m.converters {
		if err := c.Finish(); err != nil {
			return fmt.Errorf("finish on converter for mimetype %s: %v", mimetype, err)
		}
	}

	return nil
}
