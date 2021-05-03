package workers

import (
	"crypto/sha256"
	"fmt"
	"os"

	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
	"github.com/enricozb/pho/workers/pkg/lib/worker"
)

type hashWorker struct {
	db *gorm.DB
}

var _ worker.Worker = &hashWorker{}

func NewHashWorker(db *gorm.DB) *hashWorker {
	return &hashWorker{db: db}
}

func (w *hashWorker) Work(job jobs.Job) error {
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
		if err := w.db.Model(&paths.PathMetadata{}).Where("path_id", path.ID).Update("init_hash", hash[:]).Error; err != nil {
			return fmt.Errorf("update init_hash: %v", err)
		}
	}

	return nil
}

func computeHash(path string) ([sha256.Size]byte, error) {
	const chunkSize = 262144 // 2 ** 18

	f, err := os.Open(path)
	if err != nil {
		return [sha256.Size]byte{}, fmt.Errorf("open: %v", err)
	}
	defer f.Close()

	var bytes [chunkSize]byte
	numBytes, err := f.Read(bytes[:])
	if err != nil {
		return [sha256.Size]byte{}, fmt.Errorf("read: %v", err)
	}

	return sha256.Sum256(bytes[:numBytes]), nil
}
