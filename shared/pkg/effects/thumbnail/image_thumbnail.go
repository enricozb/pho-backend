package thumbnail

import (
	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
	"github.com/enricozb/pho/shared/pkg/effects/jpeg"
	"github.com/enricozb/pho/shared/pkg/lib/pool"
)

// imageThumbnailGenerator generates thumbnails for images, but not movies. The output file format is a jpeg.
type imageThumbnailGenerator struct {
	p pool.Pool
}

var (
	_ = registerThumbnailGenerator(files.ImageKind, []string{"image/jpeg", "image/png"}, newImageThumbnailGenerator)
)

func newImageThumbnailGenerator() thumbnailGenerator {
	return &imageThumbnailGenerator{pool.NewPool(32)}
}

func (c *imageThumbnailGenerator) Thumbnail(src, dst string) error {
	c.p.Go(func() error { return jpeg.Thumbnail(src, dst) })

	return nil
}

func (c *imageThumbnailGenerator) Finish() error {
	return c.p.Wait()
}
