package workers_test

import (
	"os"
	"path"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/enricozb/pho/shared/pkg/effects/daos"
	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
	"github.com/enricozb/pho/workers/pkg/effects/workers"
)

func setup(t *testing.T) (*require.Assertions, *sqlx.DB, daos.Dao, func()) {
	assert := require.New(t)
	db, cleanup := testutil.MockDB(t)
	dao := daos.NewDao(db)

	return assert, db, dao, cleanup
}

func TestJobs_ScanWorker(t *testing.T) {
	assert, db, dao, cleanup := setup(t)
	defer cleanup()

	cwd, err := os.Getwd()
	assert.NoError(err, "getwd")

	opts := jobs.ImportOptions{Paths: []string{path.Join(cwd, ".fixtures")}}

	importID := testutil.MockImportWithOptions(t, db, opts)
	jobID, err := dao.PushJob(importID, jobs.JobScan)
	assert.NoError(err, "push job")

	scanWorker := workers.NewScanWorker(dao)
	assert.NoError(scanWorker.Work(jobID))

	paths, err := dao.Paths(importID)
	assert.NoError(err, "paths")

	assert.Len(paths, 3)

	status, err := dao.GetImportStatus(importID)
	assert.Equal(jobs.StatusScan, status)
}
