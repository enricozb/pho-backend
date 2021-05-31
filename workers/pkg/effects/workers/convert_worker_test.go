package workers_test

import (
	"testing"

	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
)

func TestWorkers_ConvertWorker(t *testing.T) {
	assert, db, cleanup := setup(t)
	defer cleanup()

	importEntry, metadataJob := runScanWorker(t, db, testutil.MediaFixturesPath)
	metadataJobs, monitorJob := runMetadataWorker(t, db, metadataJob)
	dedupeJob := runMetadataWorkers(t, db, metadataJobs, monitorJob)
	convertJob := runDedupeWorker(t, db, dedupeJob)

	// before conversion, no extensions should be populated
	var count int64
	assert.NoError(db.Model(&files.File{}).Where("extension IS NULL").Count(&count).Error)
	assert.Equal(testutil.NumUniqueFilesInFixture, count)

	runConvertWorker(t, db, convertJob)

	// after conversion, all extensions should be populated
	assert.NoError(db.Model(&files.File{}).Where("extension IS NOT NULL").Count(&count).Error)
	assert.Equal(testutil.NumUniqueFilesInFixture, count)

	assertDidSetImportStatus(assert, db, importEntry.ID, jobs.ImportStatusConvert)
}
