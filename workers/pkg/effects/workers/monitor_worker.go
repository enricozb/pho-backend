package workers

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/workers/pkg/lib/worker"
)

type monitorWorker struct {
	db               *gorm.DB
	jobKindToEnqueue jobs.JobKind
}

type MonitorWorkerArgs struct {
	JobIDs []jobs.JobID
}

var _ worker.Worker = &monitorWorker{}

func NewMonitorWorker(db *gorm.DB, jobKindToEnqueue jobs.JobKind) *monitorWorker {
	return &monitorWorker{db: db, jobKindToEnqueue: jobKindToEnqueue}
}

func (w *monitorWorker) Work(job jobs.Job) error {
	importEntry := jobs.Import{}
	if err := w.db.Find(&importEntry, job.ImportID).Error; err != nil {
		return fmt.Errorf("find import: %v", err)
	}

	// get jobs to monitor
	args := MonitorWorkerArgs{}
	if err := job.GetArgs(&args); err != nil {
		return fmt.Errorf("get args: %v", err)
	}

	for {
		allDone := true
		for _, id := range args.JobIDs {
			if status, err := jobs.GetJobStatus(w.db, id); err != nil {
				return fmt.Errorf("get job status: %v", err)
			} else if status == jobs.JobStatusFailed {
				return fmt.Errorf("monitored job failed")
			} else if status != jobs.JobStatusDone {
				allDone = false
			}
		}

		if allDone {
			break
		}

		time.Sleep(1 * time.Second)
	}

	if _, err := jobs.PushJob(w.db, importEntry.ID, w.jobKindToEnqueue); err != nil {
		return fmt.Errorf("push job: %v", err)
	}

	return nil
}
