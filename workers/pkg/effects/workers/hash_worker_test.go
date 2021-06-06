package workers_test

import (
	"testing"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
	"github.com/enricozb/pho/workers/pkg/effects/workers"
)

func TestWorkers_HashWorker(t *testing.T) {
	assert, db, cleanup := setup(t)
	defer cleanup()

	importEntry, metadataJob := runScanWorker(t, db, testutil.MediaFixturesPath)
	metadataJobs, _ := runMetadataWorker(t, db, metadataJob)

	pathsToCheck, err := paths.PathsInPipeline(db, importEntry.ID)
	assert.NoError(err)
	assert.Len(pathsToCheck, int(testutil.NumFilesInFixture))
	for _, path := range pathsToCheck {
		assert.Len(path.InitHash, 0)
	}

	assert.NoError(workers.NewHashWorker(db).Work(metadataJobs[jobs.JobMetadataHash]))

	pathsToCheck, err = paths.PathsInPipeline(db, importEntry.ID)
	assert.NoError(err)
	assert.Len(pathsToCheck, int(testutil.NumFilesInFixture))
	for _, path := range pathsToCheck {
		assert.NotEqual(0, len(path.InitHash))
	}
}
