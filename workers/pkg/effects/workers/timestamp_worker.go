package workers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
	"github.com/enricozb/pho/workers/pkg/lib/worker"
)

type timestampWorker struct {
	db *gorm.DB
}

var _ worker.Worker = &timestampWorker{}

func NewTimestampWorker(db *gorm.DB) *timestampWorker {
	return &timestampWorker{db: db}
}

func (w *timestampWorker) Work(job jobs.Job) error {
	importEntry := jobs.Import{}
	if err := w.db.Find(&importEntry, job.ImportID).Error; err != nil {
		return fmt.Errorf("find import: %v", err)
	}

	var pathsToTimestamp []paths.Path
	if err := w.db.Where("import_id = ?", importEntry.ID).Find(&pathsToTimestamp).Error; err != nil {
		return fmt.Errorf("get paths: %v", err)
	}

	for _, path := range pathsToTimestamp {
		timestamp, err := computeTimestamp(path.Path)
		if err != nil {
			return fmt.Errorf("compute timestamp (%s): %v", path.Path, err)
		}
		if err := w.db.Model(&paths.PathMetadata{}).Where("path_id", path.ID).Update("timestamp", timestamp).Error; err != nil {
			return fmt.Errorf("update timestamp: %v", err)
		}
	}

	return nil
}

func computeTimestamp(path string) (int64, error) {
	cmd := exec.Command(
		"exiftool", path,
		"-api", "TimeZone=UTC",
		"-createdate",       // get the date
		"-dateFormat", "%s", // output nanoseconds since epoch
		"-s", "-s", "-s", // very short output
	)

	stdout, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("stdout pipe: %v", err)
	}

	// fallback to stat timestamp
	if len(stdout) == 0 {
		fi, err := os.Stat(path)
		if err != nil {
			return 0, fmt.Errorf("stat: %v", err)
		}

		return fi.ModTime().UnixNano(), nil
	}

	return strconv.ParseInt(string(bytes.TrimSpace(stdout)), 10, 64)
}
