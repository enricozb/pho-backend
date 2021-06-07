// Package thumbnail provides a ThumbnailGenerator that can make JPEG thumbnails for the output mime types of the converter package.
package thumbnail

import (
	"fmt"

	"github.com/enricozb/pho/shared/pkg/effects/converter"
)

// SupportedMimeTypes is the list of all mimetypes that can be converted.
var SupportedMimeTypes []string

// ThumbnailGenerator can make thumbnails for any mimetype in `SupportedMimeTypes`.
type ThumbnailGenerator struct {
	generators map[string]thumbnailGenerator
}

// thumbnailGenerator describes the behavior of all thumbnail generators (HEIC, quicktime, etc).
type thumbnailGenerator interface {
	// Thumbnail creates a thumbnail of `src` in `dst`, potentially non-blocking.
	// `dst` should not have any extensions, they will be added by the generator and returned.
	Thumbnail(src, dst string) (string, error)

	// Complete any remaining thumbnail generation tasks, blocking until all are done.
	Finish() error
}

var registeredGenerators = make(map[string]func() thumbnailGenerator)

// registerThumbnailer registers a thumbnail generator for a specific mimetype.
func registerThumbnailer(mimetype string, c func() thumbnailGenerator) struct{} {
	if _, alreadyRegistered := registeredGenerators[mimetype]; alreadyRegistered {
		panic(fmt.Errorf("converter already exists for mimetype %s", mimetype))
	}
	registeredGenerators[mimetype] = c

	SupportedMimeTypes = append(SupportedMimeTypes, mimetype)

	return struct{}{}
}

func NewThumbnailGenerator() *ThumbnailGenerator {
	m := &ThumbnailGenerator{generators: make(map[string]thumbnailGenerator)}
	for mimetype, c := range registeredGenerators {
		m.generators[mimetype] = c()
	}

	return m
}

func (t *ThumbnailGenerator) Thumbnail(src, dst, srcMimeType string) (string, error) {
	g, exists := t.generators[srcMimeType]
	if !exists {
		return "", fmt.Errorf("mimetype not supported: %s", srcMimeType)
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
	thumbnailInputMimeTypes := map[string]struct{}{}

	for _, mimetype := range SupportedMimeTypes {
		thumbnailInputMimeTypes[mimetype] = struct{}{}
	}

	for mimetype := range converter.OutputMimeTypes {
		if _, exists := thumbnailInputMimeTypes[mimetype]; !exists {
			panic(fmt.Errorf("converter outputs mimetype '%s', which is not supported by the thumbnail generator", mimetype))
		}
	}
}
