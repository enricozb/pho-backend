package workers_test

import (
	"testing"

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

	runConvertWorker(t, db, convertJob)

	assertDidSetImportStatus(assert, db, importEntry.ID, jobs.ImportStatusConvert)
}
