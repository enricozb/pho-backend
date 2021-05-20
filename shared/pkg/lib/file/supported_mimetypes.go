package file

import (
	"github.com/enricozb/pho/shared/pkg/effects/converter"
	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
)

// SupportedMimeTypes is a mapping of all supported mimetypes to their FileKind.

var SupportedMimeTypeKinds = map[string]files.FileKind{}

func init() {
	var mimetypeKinds = map[string]files.FileKind{
		"image/png":  files.ImageKind,
		"image/jpeg": files.ImageKind,
		"image/heic": files.ImageKind,
	}

	for _, mimetype := range converter.SupportedMimeTypes {
		SupportedMimeTypeKinds[mimetype] = mimetypeKinds[mimetype]
	}
}
