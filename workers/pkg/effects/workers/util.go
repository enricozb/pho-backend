package workers

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
)

func getJobInfo(dao jobs.Dao, jobID jobs.JobID) (jobs.ImportID, jobs.ImportOptions, error) {
	importID, err := dao.GetJobImportID(jobID)
	if err != nil {
		return uuid.Nil, jobs.ImportOptions{}, fmt.Errorf("get job import id: %v", err)
	}

	opts, err := dao.GetImportOptions(importID)
	if err != nil {
		return uuid.Nil, jobs.ImportOptions{}, fmt.Errorf("get import options: %v", err)
	}

	return importID, opts, nil
}
