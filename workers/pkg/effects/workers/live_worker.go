package workers

import (
	"bytes"
	"fmt"
	"os/exec"

	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
	"github.com/enricozb/pho/workers/pkg/lib/worker"
)

type liveWorker struct {
	db *gorm.DB
}

var _ worker.Worker = &liveWorker{}

// NewLiveWorker constructs a liveWorker, which is used to extract UUIDs from iOS "live" photos, in order to group the videos and images together.
func NewLiveWorker(db *gorm.DB) *liveWorker {
	return &liveWorker{db: db}
}

func (w *liveWorker) Work(job jobs.Job) error {
	importEntry := jobs.Import{}
	if err := w.db.Find(&importEntry, job.ImportID).Error; err != nil {
		return fmt.Errorf("find import: %v", err)
	}

	var pathsToHash []paths.Path
	if err := w.db.Where("import_id = ?", importEntry.ID).Find(&pathsToHash).Error; err != nil {
		return fmt.Errorf("get paths: %v", err)
	}

	for _, path := range pathsToHash {
		hash, err := computeHash(path.Path)
		if err != nil {
			return fmt.Errorf("compute hash: %v", err)
		}
		if err := w.db.Model(&paths.Path{}).Where("id = ?", path.ID).Update("init_hash", hash[:]).Error; err != nil {
			return fmt.Errorf("update init_hash: %v", err)
		}
	}

	return nil
}

func getLiveUUID(path string) ([]byte, error) {
	cmd := exec.Command(
		"exiftool", path,
		"-ContentIdentifier",
		"-MediaGroupUUID",
		"-s", "-s", "-s", // very short output
	)

	stdout, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("run cmd: %v", err)
	}

	if len(stdout) == 0 {
		return nil, nil
	}

	uuids := bytes.Split(bytes.TrimSpace(stdout), []byte("\n"))
	if len(uuids) > 1 {
		fmt.Printf("more than one live uuid for '%s'\n", path)
	}

	return uuids[0], nil
}
