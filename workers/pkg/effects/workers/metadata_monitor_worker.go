package workers

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/workers/pkg/lib/worker"
)

type metadataMonitorWorker struct {
	db *gorm.DB
}

type MetadataMonitorWorkerArgs struct {
	MetadataJobIDs []jobs.JobID
}

var _ worker.Worker = &metadataMonitorWorker{}

func NewMetadataMonitorWorker(db *gorm.DB) *metadataMonitorWorker {
	return &metadataMonitorWorker{db: db}
}

func (w *metadataMonitorWorker) Work(job jobs.Job) error {
	importEntry := jobs.Import{}
	if err := w.db.Find(&importEntry, job.ImportID).Error; err != nil {
		return fmt.Errorf("find import: %v", err)
	}

	// get jobs to monitor
	args := MetadataMonitorWorkerArgs{}
	if err := job.GetArgs(&args); err != nil {
		return fmt.Errorf("get args: %v", err)
	}

	for {
		allDone := true
		for _, id := range args.MetadataJobIDs {
			if status, err := jobs.GetJobStatus(w.db, id); err != nil {
				return fmt.Errorf("get job status: %v", err)
			} else if status == jobs.JobStatusFailed {
				return fmt.Errorf("metadata job failed")
			} else if status != jobs.JobStatusDone {
				allDone = false
			}
		}

		if allDone {
			break
		}

		time.Sleep(1 * time.Second)
	}

	if _, err := jobs.PushJob(w.db, importEntry.ID, jobs.JobDedupe); err != nil {
		return fmt.Errorf("push job: %v", err)
	}

	return nil
}
