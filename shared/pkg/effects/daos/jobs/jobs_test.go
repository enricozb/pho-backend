package jobs_test

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
)

func setup(t *testing.T) (*require.Assertions, *sqlx.DB, jobs.Dao, func()) {
	assert := require.New(t)
	db, cleanup := testutil.MockDB(t)
	dao := jobs.NewDao(db)

	return assert, db, dao, cleanup
}

func TestJobs_NewImport(t *testing.T) {
	assert, db, dao, cleanup := setup(t)
	defer cleanup()

	assert.Equal(0, testutil.NumRows(t, db, "imports"))

	importID, err := dao.NewImport(jobs.ImportOptions{})
	assert.NoError(err, "new import")

	assert.Equal(1, testutil.NumRows(t, db, "imports"))

	status, err := dao.GetImportStatus(importID)
	assert.NoError(err, "get import status")

	assert.Equal(jobs.StatusNotStarted, status)
}

func TestJobs_SetImportStatus(t *testing.T) {
	assert, _, dao, cleanup := setup(t)
	defer cleanup()

	importID, err := dao.NewImport(jobs.ImportOptions{})
	assert.NoError(err, "new import")

	err = dao.SetImportStatus(importID, jobs.StatusDedupe)
	assert.NoError(err, "set import status")

	status, err := dao.GetImportStatus(importID)
	assert.NoError(err, "get import status")

	assert.Equal(jobs.StatusDedupe, status)
}

func TestJobs_PushPeekPopJob(t *testing.T) {
	assert, db, dao, cleanup := setup(t)
	defer cleanup()

	importID, err := dao.NewImport(jobs.ImportOptions{})
	assert.NoError(err, "new import")

	assert.Equal(0, testutil.NumRows(t, db, "jobs"))

	jobID, err := dao.PushJob(importID, jobs.JobScan)
	assert.NoError(err, "push job")

	assert.Equal(1, testutil.NumRows(t, db, "jobs"))

	job, err := dao.PeekJob(importID)
	assert.NoError(err, "peek job")

	// should not change due to a peek
	assert.Equal(1, testutil.NumRows(t, db, "jobs"))

	assert.Equal(jobID, job.ID)
	assert.Equal(importID, job.ImportID)
	assert.Equal(jobs.JobScan, job.Kind)

	job, err = dao.PopJob(importID)
	assert.NoError(err, "pop job")

	// should change due to a pop
	assert.Equal(0, testutil.NumRows(t, db, "jobs"))

	assert.Equal(jobID, job.ID)
	assert.Equal(importID, job.ImportID)
	assert.Equal(jobs.JobScan, job.Kind)
}

func TestJobs_DeleteJob(t *testing.T) {
	assert, db, dao, cleanup := setup(t)
	defer cleanup()

	importID, err := dao.NewImport(jobs.ImportOptions{})
	assert.NoError(err, "new import")

	assert.Equal(0, testutil.NumRows(t, db, "jobs"))

	// insert some jobs...
	numJobs := 10
	jobIDs := make([]jobs.JobID, numJobs)
	for i := range jobIDs {
		jobIDs[i], err = dao.PushJob(importID, jobs.JobScan)
		assert.NoError(err, "push job")
	}

	assert.Equal(numJobs, testutil.NumRows(t, db, "jobs"))

	for _, jobID := range jobIDs {
		assert.NoError(dao.DeleteJob(jobID), "delete job")
		numJobs -= 1

		allJobs, err := dao.AllJobs(importID)
		assert.NoError(err, "all jobs")

		assert.Equal(numJobs, len(allJobs))

		// none of the remaining jobs should have the same ID as the deleted job
		for _, job := range allJobs {
			assert.NotEqual(jobID, job.ID)
		}
	}
}
