package workers_test

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/enricozb/pho/shared/pkg/effects/daos"
	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
)

func setup(t *testing.T) (*require.Assertions, *sqlx.DB, daos.Dao, func()) {
	assert := require.New(t)
	db, cleanup := testutil.MockDB(t)
	dao := daos.NewDao(db)

	return assert, db, dao, cleanup
}

func assertDidSetImportStatus(assert *require.Assertions, dao jobs.Dao, importID jobs.ImportID, expectedStatus jobs.Status) {
	actualStatus, err := dao.GetImportStatus(importID)
	assert.NoError(err, "get import status")

	assert.Equal(expectedStatus, actualStatus)
}

func assertDidEnqueueJob(assert *require.Assertions, dao jobs.Dao, importID jobs.ImportID, expectedKind jobs.JobKind) {
	jobs, err := dao.AllJobs()
	assert.NoError(err, "all jobs")

	found := false

	for _, job := range jobs {
		if job.ImportID == importID && job.Kind == expectedKind {
			found = true
			break
		}
	}

	assert.True(found, "found matching job kind")
}
