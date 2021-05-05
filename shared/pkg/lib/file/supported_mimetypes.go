package file

import "github.com/enricozb/pho/shared/pkg/effects/daos/files"

// SupportedMimeTypes is a mapping of all supported mimetypes to their FileKind.

var SupportedMimeTypes = map[string]files.FileKind{
	"image/png":  files.ImageKind,
	"image/jpeg": files.ImageKind,
	"image/heic": files.ImageKind,

	"video/quicktime": files.VideoKind,
}
