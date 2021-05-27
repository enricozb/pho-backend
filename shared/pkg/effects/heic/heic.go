package heic

import (
	"fmt"
	"os/exec"
)

// Convert converts src (HEIC) to dst (JPG or PNG) using a quality of 92.
func Convert(src, dst string) error {
	output, err := exec.Command("heif-convert", src, dst, "-q", "92").Output()
	if err != nil {
		return fmt.Errorf("heif-convert (%s -> %s): %v\nstderr: %s", src, dst, err, output)
	}

	return nil
}
