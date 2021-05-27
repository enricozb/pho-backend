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

	var count int64
	assert.NoError(db.Model(&paths.Path{}).Where("import_id = ?", importEntry.ID).Count(&count).Error)
	assert.Equal(testutil.NumFilesInFixture, count)

	assert.NoError(db.Model(&paths.Path{}).Where("exif_metadata = x'7B7D'").Count(&count).Error)
	assert.Equal(testutil.NumFilesInFixture, count)

	assert.NoError(workers.NewEXIFWorker(db).Work(exifJob))

	assert.NoError(db.Model(&paths.Path{}).Where("exif_metadata != x'7B7D'").Count(&count).Error)
	assert.Equal(testutil.NumFilesInFixture, count)
}
