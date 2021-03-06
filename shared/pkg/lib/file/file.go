package file

import (
	"os"

	"github.com/gabriel-vasile/mimetype"
)

// IsDir returns whether or not the path is a directory. If any error occurs, IsDir returns false.
func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

// IsSupported returns whether or not the file is of a supported format. If any error occurs, IsSupported returns false.
func IsSupported(path string) bool {
	mime, err := mimetype.DetectFile(path)
	if err != nil {
		return false
	}

	_, ok := SupportedMimeTypes[mime.String()]
	return ok
}
