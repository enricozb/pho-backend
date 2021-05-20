package converter

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/enricozb/pho/shared/pkg/effects/copyfile"
)

type identityConverter struct {
	ctx context.Context
	g   *errgroup.Group
}

func init() {
	registerConverter("image/png", newHEICConverter)
	registerConverter("image/jpeg", newHEICConverter)
}

func newIdentityConverter() Converter {
	g, ctx := errgroup.WithContext(context.Background())
	return &identityConverter{ctx: ctx, g: g}
}

func (c *identityConverter) Convert(src, dst string) error {
	c.g.Go(func() error { return copyfile.CopyFile(src, dst) })

	return nil
}

func (c *identityConverter) Finish() error {
	return c.g.Wait()
}
