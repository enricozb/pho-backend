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

type CleanupWorkerArgs struct {
	// If Full, delete all paths for an import. If not Full, do not delete failed paths.
	Full bool
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

	args := CleanupWorkerArgs{}
	if err := job.GetArgs(&args); err != nil {
		return fmt.Errorf("get args: %v", err)
	}

	if err := importEntry.SetStatus(w.db, jobs.ImportStatusCleanup); err != nil {
		return fmt.Errorf("set import status cleanup: %v", err)
	}

	if args.Full {
		if err := w.db.Where("import_id = ?", importEntry.ID).Delete(&paths.Path{}).Error; err != nil {
			return fmt.Errorf("delete all paths: %v", err)
		}
	} else {
		pathsToDelete, err := paths.PathsInPipeline(w.db, importEntry.ID)
		if err != nil {
			return fmt.Errorf("get valid paths: %v", err)
		}

		for _, path := range pathsToDelete {
			if err := w.db.Delete(path).Error; err != nil {
				return fmt.Errorf("delete path: %v", err)
			}
		}
	}

	if err := importEntry.SetStatus(w.db, jobs.ImportStatusDone); err != nil {
		return fmt.Errorf("set import status done: %v", err)
	}

	return nil
}
