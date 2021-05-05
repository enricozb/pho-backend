package workers_test

import (
	"testing"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
	"github.com/enricozb/pho/workers/pkg/effects/workers"
)

func TestWorkers_HashWorker(t *testing.T) {
	assert, db, cleanup := setup(t)
	defer cleanup()

	importEntry, metadataJob := runScanWorker(t, db, ".fixtures")
	metadataJobs, _ := runMetadataWorker(t, db, metadataJob)

	var count int64
	assert.NoError(db.Model(&paths.Path{}).Where("import_id = ?", importEntry.ID).Count(&count).Error)
	assert.Equal(numFilesInFixture, count)

	assert.NoError(db.Model(&paths.Path{}).Where("init_hash IS NULL").Count(&count).Error)
	assert.Equal(numFilesInFixture, count)

	assert.NoError(workers.NewHashWorker(db).Work(metadataJobs[jobs.JobMetadataHash]))

	assert.NoError(db.Model(&paths.Path{}).Where("init_hash IS NOT NULL").Count(&count).Error)
	assert.Equal(numFilesInFixture, count)

}
