package workers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/config"
	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
	"github.com/enricozb/pho/workers/pkg/effects/workers"
)

func setup(t *testing.T) (*require.Assertions, *gorm.DB, func()) {
	assert := require.New(t)
	db, cleanup := testutil.MockDB(t)

	tmp, err := os.MkdirTemp("", "pho-tests-datapath-*")
	assert.NoError(err)

	// modify the config so no user data will be polluted
	config.Config.DataPath = tmp

	return assert, db, func() {
		cleanup()
		os.RemoveAll(tmp)
	}
}

func assertDidSetImportStatus(assert *require.Assertions, db *gorm.DB, importID jobs.ImportID, expectedStatus jobs.ImportStatus) {
	importEntry := jobs.Import{ID: importID}
	assert.NoError(db.First(&importEntry).Error)
	assert.Equal(expectedStatus, importEntry.Status)
}

func assertDidEnqueueJob(assert *require.Assertions, db *gorm.DB, importID jobs.ImportID, kind jobs.JobKind) {
	var count int64
	assert.NoError(db.Model(&jobs.Job{}).Where("import_id = ? AND kind = ?", importID, kind).Count(&count).Error)
	assert.Equal(int64(1), count)
}

func assertDidNotEnqueueJob(assert *require.Assertions, db *gorm.DB, importID jobs.ImportID, kind jobs.JobKind) {
	var count int64
	assert.NoError(db.Model(&jobs.Job{}).Where("import_id = ? AND kind = ?", importID, kind).Count(&count).Error)
	assert.Equal(int64(0), count)
}

func runScanWorker(t *testing.T, db *gorm.DB, inputPath string) (importEntry jobs.Import, metadataJob jobs.Job) {
	assert := require.New(t)

	assert.True(filepath.IsAbs(inputPath))

	importEntry = testutil.MockImportWithOptions(t, db, jobs.ImportOptions{Paths: []string{inputPath}})
	job, err := jobs.PushJob(db, importEntry.ID, jobs.JobScan)
	assert.NoError(err, "push job")

	assert.NoError(workers.NewScanWorker(db).Work(job))

	assertDidEnqueueJob(assert, db, importEntry.ID, jobs.JobMetadata)
	assert.NoError(db.Where("import_id = ? AND kind = ?", importEntry.ID, jobs.JobMetadata).First(&metadataJob).Error)

	return importEntry, metadataJob
}

func runMetadataWorker(t *testing.T, db *gorm.DB, metadataJob jobs.Job) (metadataJobs map[jobs.JobKind]jobs.Job, monitorJob jobs.Job) {
	assert := require.New(t)

	assert.NoError(workers.NewMetadataWorker(db).Work(metadataJob))

	metadataJobs = map[jobs.JobKind]jobs.Job{}
	for _, kind := range jobs.MetadataJobKinds {
		var job jobs.Job
		assert.NoError(db.Where("import_id = ? AND kind = ?", metadataJob.ImportID, kind).First(&job).Error)
		metadataJobs[kind] = job
	}

	assert.NoError(db.Where("import_id = ? AND kind = ?", metadataJob.ImportID, jobs.JobMetadataMonitor).First(&monitorJob).Error)

	return metadataJobs, monitorJob
}

func runMetadataWorkers(t *testing.T, db *gorm.DB, metadataJobs map[jobs.JobKind]jobs.Job, monitorJob jobs.Job) (dedupeJob jobs.Job) {
	assert := require.New(t)

	for jobKind, job := range metadataJobs {
		switch jobKind {
		case jobs.JobMetadataHash:
			runHashWorker(t, db, job)
		case jobs.JobMetadataEXIF:
			runEXIFWorker(t, db, job)
		default:
			assert.Fail("unexpected job kind: " + string(jobKind))
		}

		assert.NoError(job.SetStatus(db, jobs.JobStatusDone))
	}

	assert.NoError(workers.NewMonitorWorker(db, jobs.JobDedupe).Work(monitorJob))
	assert.NoError(db.Where("import_id = ? AND kind = ?", monitorJob.ImportID, jobs.JobDedupe).First(&dedupeJob).Error)
	return dedupeJob
}

func runHashWorker(t *testing.T, db *gorm.DB, hashJob jobs.Job) {
	assert := require.New(t)
	assert.NoError(workers.NewHashWorker(db).Work(hashJob))
}

func runEXIFWorker(t *testing.T, db *gorm.DB, exifJob jobs.Job) {
	assert := require.New(t)
	assert.NoError(workers.NewEXIFWorker(db).Work(exifJob))
}

func runDedupeWorker(t *testing.T, db *gorm.DB, dedupeJob jobs.Job) (convertJob jobs.Job) {
	assert := require.New(t)
	assert.NoError(workers.NewDedupeWorker(db).Work(dedupeJob))

	assert.NoError(db.Where("import_id = ? AND kind = ?", dedupeJob.ImportID, jobs.JobConvert).First(&convertJob).Error)

	return convertJob
}

func runConvertWorker(t *testing.T, db *gorm.DB, convertJob jobs.Job) (cleanupJob jobs.Job) {
	assert := require.New(t)
	assert.NoError(workers.NewConvertWorker(db).Work(convertJob))

	assert.NoError(db.Where("import_id = ? AND kind = ?", convertJob.ImportID, jobs.JobCleanup).First(&cleanupJob).Error)

	return cleanupJob
}
