package workers

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/workers/pkg/lib/worker"
)

type metadataWorker struct {
	db *gorm.DB
}

var _ worker.Worker = &metadataWorker{}

func NewMetadataWorker(db *gorm.DB) *metadataWorker {
	return &metadataWorker{db: db}
}

func (w *metadataWorker) Work(job jobs.Job) error {
	importEntry := jobs.Import{}
	if err := w.db.Find(&importEntry, job.ImportID).Error; err != nil {
		return fmt.Errorf("find import: %v", err)
	}

	importEntry.Status = jobs.ImportStatusMetadata
	if err := w.db.Save(&importEntry).Error; err != nil {
		return fmt.Errorf("set import status: %v", err)
	}

	for _, metadataJobKind := range jobs.MetadataJobKinds {
		if _, err := jobs.PushJob(w.db, importEntry.ID, metadataJobKind); err != nil {
			return fmt.Errorf("push job (%s): %v", metadataJobKind, err)
		}
	}

	if _, err := jobs.PushJob(w.db, importEntry.ID, jobs.JobMetadataMonitor); err != nil {
		return fmt.Errorf("push monitor job: %v", err)
	}

	return nil
}
