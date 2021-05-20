package converter

import (
	"context"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/enricozb/pho/shared/pkg/effects/copyfile"
)

// identityConverter copies from src to dst, and does no conversion between media formats.
type identityConverter struct {
	ctx context.Context
	g   *errgroup.Group
}

func init() {
	registerConverter("image/png", newHEICConverter)
	registerConverter("image/jpeg", newHEICConverter)
}

func newIdentityConverter() converter {
	g, ctx := errgroup.WithContext(context.Background())
	return &identityConverter{ctx: ctx, g: g}
}

func (c *identityConverter) Convert(src, dst string) error {
	dst = dst + strings.ToUpper(filepath.Ext(src))
	c.g.Go(func() error { return copyfile.CopyFile(src, dst) })

	return nil
}

func (c *identityConverter) Finish() error {
	return c.g.Wait()
}
