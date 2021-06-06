package workers_test

import (
	"testing"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
	"github.com/enricozb/pho/workers/pkg/effects/workers"
)

func TestWorkers_EXIFWorker(t *testing.T) {
	assert, db, cleanup := setup(t)
	defer cleanup()

	importEntry, _ := runScanWorker(t, db, testutil.MediaFixturesPath)
	exifJob, err := jobs.PushJob(db, importEntry.ID, jobs.JobMetadataEXIF)
	assert.NoError(err)

	pathsToCheck, err := paths.PathsInPipeline(db, importEntry.ID)
	assert.NoError(err)
	assert.Len(pathsToCheck, int(testutil.NumFilesInFixture))

	// check that exif metadata is empty for valid paths
	for _, path := range pathsToCheck {
		assert.Equal(paths.EXIFMetadata{}, path.EXIFMetadata)
	}

	assert.NoError(workers.NewEXIFWorker(db).Work(exifJob))

	pathsToCheck, err = paths.PathsInPipeline(db, importEntry.ID)
	assert.NoError(err)
	assert.Len(pathsToCheck, int(testutil.NumFilesInFixture))

	// check that exif metadata is not empty for valid paths
	for _, path := range pathsToCheck {
		assert.NotEqual(path.EXIFMetadata, paths.EXIFMetadata{})
	}
}
