package thumbnail

import (
	"github.com/enricozb/pho/shared/pkg/effects/jpeg"
	"github.com/enricozb/pho/shared/pkg/lib/pool"
)

// imageThumbnailGenerator generates thumbnails for images, but not movies. The output file format is a jpeg.
type imageThumbnailGenerator struct {
	p pool.Pool
}

var (
	_ = registerThumbnailer("image/jpeg", newImageThumbnailGenerator)
	_ = registerThumbnailer("image/png", newImageThumbnailGenerator)
)

func newImageThumbnailGenerator() thumbnailGenerator {
	return &imageThumbnailGenerator{pool.NewPool(32)}
}

func (c *imageThumbnailGenerator) Thumbnail(src, dst string) (string, error) {
	dst = dst + ".JPG"
	c.p.Go(func() error { return jpeg.Thumbnail(src, dst) })

	return dst, nil
}

func (c *imageThumbnailGenerator) Finish() error {
	return c.p.Wait()
}
