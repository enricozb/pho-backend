package jpeg

import (
	"fmt"
	"os/exec"
)

// Thumbnail creates a thumbnail of a `src` file (any image mimetype) with `dst` as the destination.
// The aspect ratio is preserved.
func Thumbnail(src, dst string) error {
	// TODO(enricozb): add some check somewhere that if the size is less than X, just copy the file.

	output, err := exec.Command(
		"convert",
		"-define", "jpeg:size=480x360",
		"-auto-orient",
		"-thumbnail", "120x90",
		"-strip",
		src,
		dst,
	).Output()

	if err != nil {
		return fmt.Errorf("thumbnail (%s -> %s): %v\nstderr: %s", src, dst, err, output)
	}

	return nil
}
