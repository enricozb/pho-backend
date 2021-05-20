package converter

import (
	"github.com/enricozb/pho/shared/pkg/effects/heic"
	"github.com/enricozb/pho/shared/pkg/lib/pool"
)

type heicConverter struct {
	p pool.Pool
}

var _ Converter = &heicConverter{}

func newHEICConverter() *heicConverter {
	return &heicConverter{pool.NewPool(32)}
}

func (c *heicConverter) Convert(src, dst string) error {
	c.p.Go(func() error { return heic.Convert(src, dst) })

	return nil
}

func (c *heicConverter) Finish() error {
	return c.p.Wait()
}
