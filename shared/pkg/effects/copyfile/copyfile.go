package copyfile

import (
	"fmt"
	"io"
	"os"
	"syscall"
)

// CopyFile copies files with persmissions, ownerships, and times.
func CopyFile(src, dst string) error {
	stat, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat: %v", err)
	}

	srcf, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open: %v", err)
	}
	defer srcf.Close()

	dstf, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create: %v", err)
	}
	defer dstf.Close()

	if _, err := io.Copy(dstf, srcf); err != nil {
		return fmt.Errorf("copy: %v", err)
	}

	if err := dstf.Chmod(stat.Mode()); err != nil {
		return fmt.Errorf("chmod: %v", err)
	}

	if unixstat, ok := stat.Sys().(*syscall.Stat_t); !ok {
		if err := dstf.Chown(int(unixstat.Uid), int(unixstat.Gid)); err != nil {
			return fmt.Errorf("chown: %v", err)
		}
		return fmt.Errorf("stat to Stat_t")
	}

	if err := os.Chtimes(dst, stat.ModTime(), stat.ModTime()); err != nil {
		return fmt.Errorf("chtimes: %v", err)
	}

	return nil
}
