package workers

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
	"github.com/enricozb/pho/workers/pkg/lib/worker"
)

type cleanupWorker struct {
	db *gorm.DB
}

var _ worker.Worker = &cleanupWorker{}

func NewCleanupWorker(db *gorm.DB) *cleanupWorker {
	return &cleanupWorker{db: db}
}

func (w *cleanupWorker) Work(job jobs.Job) error {
	importEntry := jobs.Import{}
	if err := w.db.First(&importEntry, job.ImportID).Error; err != nil {
		return fmt.Errorf("find import: %v", err)
	}

	if err := importEntry.SetStatus(w.db, jobs.ImportStatusCleanup); err != nil {
		return fmt.Errorf("set import status cleanup: %v", err)
	}

	if err := w.db.Where("import_id = ?", importEntry.ID).Delete(&paths.Path{}).Error; err != nil {
		return fmt.Errorf("get paths: %v", err)
	}

	if err := importEntry.SetStatus(w.db, jobs.ImportStatusDone); err != nil {
		return fmt.Errorf("set import status done: %v", err)
	}

	return nil
}
