package workers_test

import (
	"os"
	"path"
	"testing"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
	"github.com/enricozb/pho/workers/pkg/effects/workers"
)

func TestWorkers_ScanWorker(t *testing.T) {
	assert, db, cleanup := setup(t)
	defer cleanup()

	cwd, err := os.Getwd()
	assert.NoError(err, "getwd")

	opts := jobs.ImportOptions{Paths: []string{path.Join(cwd, ".fixtures")}}

	importEntry := testutil.MockImportWithOptions(t, db, opts)
	job, err := jobs.PushJob(db, importEntry.ID, jobs.JobScan)
	assert.NoError(err, "push job")

	scanWorker := workers.NewScanWorker(db)
	assert.NoError(scanWorker.Work(job))

	var count int64
	assert.NoError(db.Model(&paths.Path{}).Where("import_id = ?", importEntry.ID).Count(&count).Error)

	assert.Equal(int64(4), count)

	assertDidSetImportStatus(assert, db, importEntry.ID, jobs.ImportStatusScan)
	assertDidEnqueueJob(assert, db, importEntry.ID, jobs.JobMetadata)
}
