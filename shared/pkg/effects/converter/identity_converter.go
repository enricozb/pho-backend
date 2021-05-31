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
	registerConverter("image/png", newIdentityConverter)
	registerConverter("image/jpeg", newIdentityConverter)
}

func newIdentityConverter() converter {
	g, ctx := errgroup.WithContext(context.Background())
	return &identityConverter{ctx: ctx, g: g}
}

func (c *identityConverter) Convert(src, dst string) (string, error) {
	dst = dst + strings.ToUpper(filepath.Ext(src))
	c.g.Go(func() error { return copyfile.CopyFile(src, dst) })

	return dst, nil
}

func (c *identityConverter) Finish() error {
	return c.g.Wait()
}
