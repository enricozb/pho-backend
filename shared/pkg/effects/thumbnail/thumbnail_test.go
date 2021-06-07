package thumbnail_test

import (
	"testing"

	"github.com/enricozb/pho/shared/pkg/effects/thumbnail"
)

// Test_Thumbnail_SupportedMimeTypes tests that the thumbnail generator can take in as input all possible output formats for the converter.
// The thumbnail package has an `init` function that tests this, so importing it alone will test it.
func Test_Thumbnail_SupportedMimeTypes(t *testing.T) {
	_ = thumbnail.NewThumbnailGenerator()
}
