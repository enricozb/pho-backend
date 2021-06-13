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

	// heif-convert images need to have their output orientations set to 1
	// see: https://github.com/enricozb/pho-backend/issues/3
	output, err = exec.Command("exiftool", dst, "-orientation#=1").Output()
	if err != nil {
		return fmt.Errorf("exiftool %s -orientation#=1: %v\nstderr: %s", dst, err, output)
	}

	return nil
}
