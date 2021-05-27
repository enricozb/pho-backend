package workers_test

import (
	"testing"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
	"github.com/enricozb/pho/workers/pkg/effects/workers"
)

func TestWorkers_ScanWorker(t *testing.T) {
	assert, db, cleanup := setup(t)
	defer cleanup()

	opts := jobs.ImportOptions{Paths: []string{testutil.MediaFixturesPath}}

	importEntry := testutil.MockImportWithOptions(t, db, opts)
	job, err := jobs.PushJob(db, importEntry.ID, jobs.JobScan)
	assert.NoError(err, "push job")

	scanWorker := workers.NewScanWorker(db)
	assert.NoError(scanWorker.Work(job))

	var count int64
	assert.NoError(db.Model(&paths.Path{}).Where("import_id = ? AND kind != ''", importEntry.ID).Count(&count).Error)
	assert.Equal(numFilesInFixture, count)

	assertDidSetImportStatus(assert, db, importEntry.ID, jobs.ImportStatusScan)
	assertDidEnqueueJob(assert, db, importEntry.ID, jobs.JobMetadata)
}

func TestWorkers_ScanWorker_DuplicatePaths(t *testing.T) {
	assert, db, cleanup := setup(t)
	defer cleanup()

	opts := jobs.ImportOptions{Paths: []string{testutil.MediaFixturesPath, testutil.MediaFixturesPath}}

	importEntry := testutil.MockImportWithOptions(t, db, opts)
	job, err := jobs.PushJob(db, importEntry.ID, jobs.JobScan)
	assert.NoError(err, "push job")

	scanWorker := workers.NewScanWorker(db)
	assert.NoError(scanWorker.Work(job))

	var count int64
	assert.NoError(db.Model(&paths.Path{}).Where("import_id = ?", importEntry.ID).Count(&count).Error)

	assert.Equal(numFilesInFixture, count)

	assertDidSetImportStatus(assert, db, importEntry.ID, jobs.ImportStatusScan)
	assertDidEnqueueJob(assert, db, importEntry.ID, jobs.JobMetadata)
}
