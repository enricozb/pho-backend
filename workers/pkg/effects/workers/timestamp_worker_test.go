package workers_test

import (
	"testing"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
	"github.com/enricozb/pho/workers/pkg/effects/workers"
)

func TestWorkers_TimestampWorker(t *testing.T) {
	assert, db, cleanup := setup(t)
	defer cleanup()

	importEntry, metadataJob := runScanWorker(t, db)
	metadataJobs, _ := runMetadataWorker(t, db, metadataJob)

	var count int64
	assert.NoError(db.Model(&paths.Path{}).Where("import_id = ?", importEntry.ID).Count(&count).Error)
	assert.Equal(int64(3), count)

	assert.NoError(db.Model(&paths.PathMetadata{}).Where("timestamp IS NULL").Count(&count).Error)
	assert.Equal(int64(3), count)

	assert.NoError(workers.NewTimestampWorker(db).Work(metadataJobs[jobs.JobMetadataTimestamp]))

	assert.NoError(db.Model(&paths.PathMetadata{}).Where("timestamp IS NOT NULL").Count(&count).Error)
	assert.Equal(int64(3), count)

}
