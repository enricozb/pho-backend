package workers_test

import (
	"testing"
	"time"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
	"github.com/enricozb/pho/workers/pkg/effects/workers"
)

const jobKindToEnqueue = jobs.JobDedupe

func TestWorkers_MetadataMonitorWorker_Enqueue(t *testing.T) {
	assert, db, cleanup := setup(t)
	defer cleanup()

	importEntry := testutil.MockImport(t, db)
	job, err := jobs.PushJobWithArgs(db, importEntry.ID, jobs.JobMetadataMonitor, workers.MonitorWorkerArgs{})
	assert.NoError(err, "push job")

	monitorWorker := workers.NewMonitorWorker(db, jobKindToEnqueue)
	assert.NoError(monitorWorker.Work(job))

	// check that no status was set
	assertDidSetImportStatus(assert, db, importEntry.ID, jobs.ImportStatusNotStarted)
	assertDidEnqueueJob(assert, db, importEntry.ID, jobKindToEnqueue)
}

func TestWorkers_MetadataMonitorWorker_Monitors(t *testing.T) {
	assert, db, cleanup := setup(t)
	defer cleanup()

	importEntry := testutil.MockImport(t, db)

	// enqueue some scan jobs to be monitored, but start none of them
	const numJobs = 3
	scanJobs := make([]jobs.Job, numJobs)
	jobIDs := make([]jobs.JobID, numJobs)
	for i := range jobIDs {
		scanJob, err := jobs.PushJob(db, importEntry.ID, jobs.JobScan)
		assert.NoError(err, "push job")

		scanJobs[i] = scanJob
		jobIDs[i] = scanJob.ID
	}

	job, err := jobs.PushJobWithArgs(db, importEntry.ID, jobs.JobMetadataMonitor, workers.MonitorWorkerArgs{JobIDs: jobIDs})
	assert.NoError(err, "push job")
	monitorWorker := workers.NewMonitorWorker(db, jobKindToEnqueue)
	go func() {
		assert.NoError(monitorWorker.Work(job))
	}()

	for _, job := range scanJobs {
		assertDidNotEnqueueJob(assert, db, importEntry.ID, jobKindToEnqueue)
		job.Status = jobs.JobStatusDone
		assert.NoError(db.Save(&job).Error)
		time.Sleep(2 * time.Second)
	}

	assertDidNotEnqueueJob(assert, db, importEntry.ID, jobKindToEnqueue)
}
