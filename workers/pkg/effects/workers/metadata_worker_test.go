package workers_test

import (
	"testing"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
	"github.com/enricozb/pho/workers/pkg/effects/workers"
)

func TestWorkers_MetadataWorker(t *testing.T) {
	assert, db, cleanup := setup(t)
	defer cleanup()

	importEntry := testutil.MockImport(t, db)
	jobID, err := jobs.PushJob(db, importEntry.ID, jobs.JobMetadata)
	assert.NoError(err, "push job")

	metadataWorker := workers.NewMetadataWorker(db)
	assert.NoError(metadataWorker.Work(jobID))

	// Check that the monitor and all metadata jobs were enqueued.
	assertDidEnqueueJob(assert, db, importEntry.ID, jobs.JobMetadataMonitor)
	for _, metadataJobKind := range jobs.MetadataJobKinds {
		assertDidEnqueueJob(assert, db, importEntry.ID, metadataJobKind)
	}

	assertDidSetImportStatus(assert, db, importEntry.ID, jobs.ImportStatusMetadata)
}
