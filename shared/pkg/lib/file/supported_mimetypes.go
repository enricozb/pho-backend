package file

// SupportedMimeTypes is a set of strings of the supported mime types for imports.
// This is exported as a map for fast lookup
var SupportedMimeTypes = map[string]interface{}{}

var supportedMimeTypes = []string{
	"image/png",
	"image/jpeg",
}

func init() {
	for _, mimetype := range supportedMimeTypes {
		SupportedMimeTypes[mimetype] = nil
	}
}
