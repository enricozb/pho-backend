package workers_test

import (
	"testing"

	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/workers/pkg/effects/workers"
)

func TestWorkers_DedupeWorker(t *testing.T) {
	assert, db, _ := setup(t)
	// defer cleanup()

	importEntry, metadataJob := runScanWorker(t, db, ".fixtures")
	metadataJobs, _ := runMetadataWorker(t, db, metadataJob)

	runHashWorker(t, db, metadataJobs[jobs.JobMetadataHash])
	runEXIFWorker(t, db, metadataJobs[jobs.JobMetadataEXIF])

	var count int64
	assert.NoError(db.Model(&files.File{}).Where("import_id = ?", importEntry.ID).Count(&count).Error)
	assert.Equal(int64(0), count)

	assert.NoError(workers.NewDedupeWorker(db).Work(metadataJobs[jobs.JobMetadataHash]))

	assert.NoError(db.Model(&files.File{}).Where("import_id = ?", importEntry.ID).Count(&count).Error)
	assert.Equal(numUniqueFilesInFixture, count)
}
