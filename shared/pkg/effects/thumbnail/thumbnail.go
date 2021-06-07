// Package thumbnail provides a ThumbnailGenerator that can make JPEG thumbnails for the output mime types of the converter package.
package thumbnail

import (
	"fmt"

	"github.com/enricozb/pho/shared/pkg/effects/converter"
	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
)

const ThumbnailSuffix string = ".JPG"

// The set of supported input mimetypes.
var SupportedMimeTypes = make(map[string]struct{})

// ThumbnailGenerator can make thumbnails for any mimetype in `SupportedMimeTypes`.
type ThumbnailGenerator struct {
	generators map[files.FileKind]thumbnailGenerator
}

// thumbnailGenerator describes the behavior of all thumbnail generators (HEIC, quicktime, etc).
type thumbnailGenerator interface {
	// Thumbnail creates a thumbnail of `src` in `dst`, potentially non-blocking.
	// The extension of `dst` must be `ThumbnailSuffix`.
	Thumbnail(src, dst string) error

	// Complete any remaining thumbnail generation tasks, blocking until all are done.
	Finish() error
}

var registeredThumbnailGenerators = make(map[files.FileKind]func() thumbnailGenerator)

// registerThumbnailGenerator registers a thumbnail generator for a specific file kind and the mimetypes it supports.
func registerThumbnailGenerator(kind files.FileKind, mimetypes []string, t func() thumbnailGenerator) struct{} {
	for _, mimetype := range mimetypes {
		if _, alreadyRegistered := SupportedMimeTypes[mimetype]; alreadyRegistered {
			panic(fmt.Errorf("thumbnail generator already exists for mimetype %s", mimetype))
		}
		SupportedMimeTypes[mimetype] = struct{}{}
	}

	registeredThumbnailGenerators[kind] = t

	return struct{}{}
}

func NewThumbnailGenerator() *ThumbnailGenerator {
	m := &ThumbnailGenerator{generators: make(map[files.FileKind]thumbnailGenerator)}
	for kind, c := range registeredThumbnailGenerators {
		m.generators[kind] = c()
	}

	return m
}

func (t *ThumbnailGenerator) Thumbnail(src, dst string, kind files.FileKind) error {
	g, exists := t.generators[kind]
	if !exists {
		return fmt.Errorf("kind not supported: %s", kind)
	}

	return g.Thumbnail(src, dst)
}

func (t *ThumbnailGenerator) Finish() error {
	for mimetype, g := range t.generators {
		if err := g.Finish(); err != nil {
			return fmt.Errorf("finish on thumbnail generator for mimetype %s: %v", mimetype, err)
		}
	}

	return nil
}

func init() {
	for mimetype := range converter.OutputMimeTypes {
		if _, exists := SupportedMimeTypes[mimetype]; !exists {
			panic(fmt.Errorf("converter outputs mimetype '%s', which is not supported by the thumbnail generator", mimetype))
		}
	}
}
