package workers_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
	"github.com/enricozb/pho/workers/pkg/effects/workers"
)

func TestWorkers_DedupeWorker(t *testing.T) {
	assert, db, cleanup := setup(t)
	defer cleanup()

	importEntry, metadataJob := runScanWorker(t, db, ".fixtures")
	metadataJobs, _ := runMetadataWorker(t, db, metadataJob)

	runHashWorker(t, db, metadataJobs[jobs.JobMetadataHash])
	runEXIFWorker(t, db, metadataJobs[jobs.JobMetadataEXIF])

	var files []files.File
	var paths []paths.Path

	assert.NoError(db.Where("import_id = ?", importEntry.ID).Find(&files).Error)
	assert.Len(files, 0)

	assert.NoError(workers.NewDedupeWorker(db).Work(metadataJobs[jobs.JobMetadataHash]))

	assert.NoError(db.Where("import_id = ?", importEntry.ID).Find(&files).Error)
	assert.Len(files, int(numUniqueFilesInFixture))

	assert.NoError(db.Where("import_id = ?", importEntry.ID).Find(&paths).Error)
	assert.Len(paths, int(numFilesInFixture))

	pathIDs := make([]uuid.UUID, len(paths))
	fileIDs := make([]uuid.UUID, len(files))

	for i, path := range paths {
		pathIDs[i] = path.ID
	}

	for i, file := range files {
		fileIDs[i] = file.ID
	}

	assert.Subset(pathIDs, fileIDs)
}
