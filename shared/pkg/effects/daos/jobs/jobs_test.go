package jobs_test

import (
	"errors"
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

	importID := testutil.MockImport(t, db)
	assert.Equal(1, testutil.NumRows(t, db, "imports"))

	status, err := dao.GetImportStatus(importID)
	assert.NoError(err, "get import status")

	assert.Equal(jobs.StatusNotStarted, status)
}

func TestJobs_SetImportStatus(t *testing.T) {
	assert, db, dao, cleanup := setup(t)
	defer cleanup()

	importID := testutil.MockImport(t, db)
	err := dao.SetImportStatus(importID, jobs.StatusDedupe)
	assert.NoError(err, "set import status")

	status, err := dao.GetImportStatus(importID)
	assert.NoError(err, "get import status")

	assert.Equal(jobs.StatusDedupe, status)
}

func TestJobs_PushPeekPopJob(t *testing.T) {
	assert, db, dao, cleanup := setup(t)
	defer cleanup()

	importID := testutil.MockImport(t, db)
	assert.Equal(0, testutil.NumRows(t, db, "jobs"))

	jobID, err := dao.PushJob(importID, jobs.JobScan)
	assert.NoError(err, "push job")

	assert.Equal(1, testutil.NumRows(t, db, "jobs"))

	job, err := dao.PeekJob()
	assert.NoError(err, "peek job")

	// should not change due to a peek
	assert.Equal(1, testutil.NumRows(t, db, "jobs"))

	assert.Equal(jobID, job.ID)
	assert.Equal(importID, job.ImportID)
	assert.Equal(jobs.JobScan, job.Kind)

	job, ok, err := dao.PopJob()
	assert.NoError(err, "pop job")
	assert.True(ok, "pop job ok")

	// should change due to a pop
	assert.Equal(0, testutil.NumRows(t, db, "jobs"))

	assert.Equal(jobID, job.ID)
	assert.Equal(importID, job.ImportID)
	assert.Equal(jobs.JobScan, job.Kind)
}

func TestJobs_NumJobs(t *testing.T) {
	assert, db, dao, cleanup := setup(t)
	defer cleanup()

	importID := testutil.MockImport(t, db)
	assert.Equal(0, testutil.NumRows(t, db, "jobs"))

	// insert some jobs...
	numJobs := 10
	jobIDs := make([]jobs.JobID, numJobs)
	for i := range jobIDs {
		var err error
		jobIDs[i], err = dao.PushJob(importID, jobs.JobScan)
		assert.NoError(err, "push job")
	}

	actualNumJobs, err := dao.NumJobs()
	assert.NoError(err, "num jobs")
	assert.Equal(numJobs, actualNumJobs, "num jobs")
}

func TestJobs_DeleteJob(t *testing.T) {
	assert, db, dao, cleanup := setup(t)
	defer cleanup()

	importID := testutil.MockImport(t, db)
	assert.Equal(0, testutil.NumRows(t, db, "jobs"))

	// insert some jobs...
	numJobs := 10
	jobIDs := make([]jobs.JobID, numJobs)
	for i := range jobIDs {
		var err error
		jobIDs[i], err = dao.PushJob(importID, jobs.JobScan)
		assert.NoError(err, "push job")
	}

	assert.Equal(numJobs, testutil.NumRows(t, db, "jobs"))

	for _, jobID := range jobIDs {
		assert.NoError(dao.DeleteJob(jobID), "delete job")
		numJobs -= 1

		allJobs, err := dao.AllJobs()
		assert.NoError(err, "all jobs")

		assert.Equal(numJobs, len(allJobs))

		// none of the remaining jobs should have the same ID as the deleted job
		for _, job := range allJobs {
			assert.NotEqual(jobID, job.ID)
		}
	}
}

func TestJobs_ImportFailures(t *testing.T) {
	assert, db, dao, cleanup := setup(t)
	defer cleanup()

	errors := []error{errors.New("error1"), errors.New("error2"), errors.New("error3")}
	expectedMessages := []string{}

	importID := testutil.MockImport(t, db)

	for _, err := range errors {
		assert.NoError(dao.RecordImportFailure(importID, err))
		expectedMessages = append(expectedMessages, err.Error())
	}

	actualMessages, err := dao.GetImportFailureMessages(importID)
	assert.NoError(err, "get import failure messages")

	assert.ElementsMatch(expectedMessages, actualMessages)
}

func TestJobs_ImportOptions(t *testing.T) {
	assert, db, dao, cleanup := setup(t)
	defer cleanup()

	expectedOpts := jobs.ImportOptions{Paths: []string{"path1", "path2", "path3"}}

	importID := testutil.MockImportWithOptions(t, db, expectedOpts)

	actualOpts, err := dao.GetImportOptions(importID)
	assert.NoError(err, "get import opts")

	assert.Equal(expectedOpts, actualOpts)
}

func TestJobs_GetJobImportID(t *testing.T) {
	assert, db, dao, cleanup := setup(t)
	defer cleanup()

	expectedImportID := testutil.MockImport(t, db)

	jobID, err := dao.PushJob(expectedImportID, jobs.JobScan)
	assert.NoError(err, "push job")

	actualImportID, err := dao.GetJobImportID(jobID)
	assert.NoError(err, "get job import id")

	assert.Equal(expectedImportID, actualImportID)
}
