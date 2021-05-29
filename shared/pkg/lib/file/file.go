package file

import (
	"fmt"
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

	kind, ok = SupportedMimeTypeKinds[mimeinfo.String()]
	return ok, kind, mimeinfo.String()
}

// MakeDirIfNotExist creates the directory passed in, including the parents if the entire path doesn't exist.
func MakeDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("mkdir: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("stat: %v", err)
	}

	return nil
}
