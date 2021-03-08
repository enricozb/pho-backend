package jobs_test

import (
	"testing"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
)

func TestJobs_JobSmokeTest(t *testing.T) {
	assert, db, cleanup := setup(t)
	defer cleanup()

	insertedImport := jobs.Import{}
	db.Create(&insertedImport)

	insertedJob := jobs.Job{ImportID: insertedImport.ID, Kind: jobs.JobScan}
	db.Create(&insertedJob)

	foundJob := jobs.Job{ImportID: insertedImport.ID}
	assert.NoError(db.First(&foundJob).Error)

	assert.Equal(jobs.JobStatusNotStarted, insertedJob.Status)
	assert.Equal(insertedJob, foundJob)

	jobs := []jobs.Job{}
	assert.NoError(db.Where("import_id = ?", insertedImport.ID.String()).Find(&jobs).Error)

	assert.Len(jobs, 1)
	assert.Equal(insertedJob, jobs[0])
}
