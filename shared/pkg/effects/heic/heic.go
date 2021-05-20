package heic

import (
	"fmt"
	"os/exec"
)

// Convert converts src (HEIC) to dst (JPG or PNG) using a quality of 90.
func Convert(src, dst string) error {
	output, err := exec.Command("heif-convert", src, dst, "-q", "90").Output()
	if err != nil {
		return fmt.Errorf("heif-convert: %v\nstderr: %s", err, output)
	}

	return nil
}
