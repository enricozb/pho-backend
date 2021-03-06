package workers

import (
	"fmt"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/workers/pkg/lib/worker"
)

type metadataWorkerDao interface {
	jobs.Dao
}

type metadataWorker struct {
	dao metadataWorkerDao
}

var _ worker.Worker = &metadataWorker{}

func NewMetadataWorker(dao metadataWorkerDao) *metadataWorker {
	return &metadataWorker{dao: dao}
}

func (w *metadataWorker) Work(jobID jobs.JobID) error {
	importID, _, err := getJobInfo(w.dao, jobID)
	if err != nil {
		return fmt.Errorf("get job info: %v", err)
	}

	if err := w.dao.SetImportStatus(importID, jobs.StatusMetadata); err != nil {
		return fmt.Errorf("set import status: %v", err)
	}

	for _, metadataJobKind := range jobs.MetadataJobKinds {
		if _, err := w.dao.PushJob(importID, metadataJobKind); err != nil {
			return fmt.Errorf("push job (%s): %v", metadataJobKind, err)
		}
	}

	return nil
}
