package workers_test

import (
	"testing"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
	"github.com/enricozb/pho/workers/pkg/effects/workers"
)

func TestWorkers_MetadataWorker(t *testing.T) {
	assert, db, dao, cleanup := setup(t)
	defer cleanup()

	importID := testutil.MockImport(t, db)
	jobID, err := dao.PushJob(importID, jobs.JobMetadata)
	assert.NoError(err, "push job")

	metadataWorker := workers.NewMetadataWorker(dao)
	assert.NoError(metadataWorker.Work(jobID))

	// Check that the monitor and all metadata jobs were enqueued.
	assertDidEnqueueJob(assert, dao, importID, jobs.JobMetadataMonitor)
	for _, metadataJobKind := range jobs.MetadataJobKinds {
		assertDidEnqueueJob(assert, dao, importID, metadataJobKind)
	}

	assertDidSetImportStatus(assert, dao, importID, jobs.ImportStatusMetadata)
}
