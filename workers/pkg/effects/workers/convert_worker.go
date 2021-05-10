package workers

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/workers/pkg/lib/worker"
)

type convertWorker struct {
	db *gorm.DB
}

var _ worker.Worker = &convertWorker{}

func NewConvertWorker(db *gorm.DB) *convertWorker {
	return &convertWorker{db: db}
}

func (w *convertWorker) Work(job jobs.Job) error {
	importEntry := jobs.Import{}
	if err := w.db.Find(&importEntry, job.ImportID).Error; err != nil {
		return fmt.Errorf("find import: %v", err)
	}

	if err := importEntry.SetStatus(w.db, jobs.ImportStatusConvert); err != nil {
		return fmt.Errorf("set import status: %v", err)
	}

	jobIDs := []jobs.JobID{}

	for _, convertJobKind := range jobs.ConvertJobKinds {
		if job, err := jobs.PushJob(w.db, importEntry.ID, convertJobKind); err != nil {
			return fmt.Errorf("push job (%s): %v", convertJobKind, err)
		} else {
			jobIDs = append(jobIDs, job.ID)
		}
	}

	if _, err := jobs.PushJobWithArgs(w.db, importEntry.ID, jobs.JobConvertMonitor, MonitorWorkerArgs{JobIDs: jobIDs}); err != nil {
		return fmt.Errorf("push monitor job: %v", err)
	}

	return nil
}
