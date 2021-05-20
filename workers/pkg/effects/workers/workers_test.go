package workers_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
	"github.com/enricozb/pho/workers/pkg/effects/workers"
)

const (
	numFilesInFixture       int64 = 6
	numUniqueFilesInFixture int64 = 5
)

func setup(t *testing.T) (*require.Assertions, *gorm.DB, func()) {
	assert := require.New(t)
	db, cleanup := testutil.MockDB(t)
	return assert, db, cleanup
}

func assertDidSetImportStatus(assert *require.Assertions, db *gorm.DB, importID jobs.ImportID, expectedStatus jobs.ImportStatus) {
	importEntry := jobs.Import{ID: importID}
	assert.NoError(db.Find(&importEntry).Error)
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

	cwd, err := os.Getwd()
	assert.NoError(err, "getwd")

	importEntry = testutil.MockImportWithOptions(t, db, jobs.ImportOptions{Paths: []string{path.Join(cwd, inputPath)}})
	job, err := jobs.PushJob(db, importEntry.ID, jobs.JobScan)
	assert.NoError(err, "push job")

	assert.NoError(workers.NewScanWorker(db).Work(job))

	assertDidEnqueueJob(assert, db, importEntry.ID, jobs.JobMetadata)
	assert.NoError(db.Where("import_id = ? AND kind = ?", importEntry.ID, jobs.JobMetadata).Find(&metadataJob).Error)

	return importEntry, metadataJob
}

func runMetadataWorker(t *testing.T, db *gorm.DB, metadataJob jobs.Job) (metadataJobs map[jobs.JobKind]jobs.Job, monitorJob jobs.Job) {
	assert := require.New(t)

	assert.NoError(workers.NewMetadataWorker(db).Work(metadataJob))

	metadataJobs = map[jobs.JobKind]jobs.Job{}
	for _, kind := range jobs.MetadataJobKinds {
		var job jobs.Job
		assert.NoError(db.Where("import_id = ? AND kind = ?", metadataJob.ImportID, kind).Find(&job).Error)
		metadataJobs[kind] = job
	}

	assert.NoError(db.Where("import_id = ? AND kind = ?", metadataJob.ImportID, jobs.JobMetadataMonitor).Find(&monitorJob).Error)

	return metadataJobs, monitorJob
}

func runHashWorker(t *testing.T, db *gorm.DB, hashJob jobs.Job) {
	assert := require.New(t)
	assert.NoError(workers.NewHashWorker(db).Work(hashJob))
}

func runEXIFWorker(t *testing.T, db *gorm.DB, exifJob jobs.Job) {
	assert := require.New(t)
	assert.NoError(workers.NewEXIFWorker(db).Work(exifJob))
}
