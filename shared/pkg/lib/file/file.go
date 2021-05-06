package file

import (
	"os"

	"github.com/gabriel-vasile/mimetype"

	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
)

// IsDir returns whether or not the path is a directory. If any error occurs, IsDir returns false.
func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

// Kind returns whether or not the file is of a supported format, the kind of the file, and the mimetype as a string. If any error occurs, `ok` is false.
func Kind(path string) (ok bool, kind files.FileKind, mime string) {
	mimeinfo, err := mimetype.DetectFile(path)
	if err != nil {
		return false, "", ""
	}

	kind, ok = SupportedMimeTypes[mimeinfo.String()]
	return ok, kind, mimeinfo.String()
}
