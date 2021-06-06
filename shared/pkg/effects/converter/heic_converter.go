package converter

import (
	"github.com/enricozb/pho/shared/pkg/effects/heic"
	"github.com/enricozb/pho/shared/pkg/lib/pool"
)

type heicConverter struct {
	p pool.Pool
}

func init() {
	registerConverter("image/heic", "image/jpeg", newHEICConverter)
}

func newHEICConverter() converter {
	return &heicConverter{pool.NewPool(32)}
}

func (c *heicConverter) Convert(src, dst string) (string, error) {
	dst = dst + ".JPG"
	c.p.Go(func() error { return heic.Convert(src, dst) })

	return dst, nil
}

func (c *heicConverter) Finish() error {
	return c.p.Wait()
}
