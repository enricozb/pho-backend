package workers_test

import (
	"os"
	"path"
	"testing"

	"gorm.io/gorm"

	"github.com/stretchr/testify/require"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
	"github.com/enricozb/pho/workers/pkg/effects/workers"
)

func runScanWorker(t *testing.T, db *gorm.DB) (importEntry jobs.Import, metadataJob jobs.Job) {
	assert := require.New(t)

	cwd, err := os.Getwd()
	assert.NoError(err, "getwd")

	importEntry = testutil.MockImportWithOptions(t, db, jobs.ImportOptions{Paths: []string{path.Join(cwd, ".fixtures")}})
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

func TestWorkers_HashWorker(t *testing.T) {
	assert, db, _ := setup(t)
	// defer cleanup()

	importEntry, metadataJob := runScanWorker(t, db)
	metadataJobs, _ := runMetadataWorker(t, db, metadataJob)

	var count int64
	assert.NoError(db.Model(&paths.Path{}).Where("import_id = ?", importEntry.ID).Count(&count).Error)
	assert.Equal(int64(3), count)

	assert.NoError(db.Model(&paths.PathMetadata{}).Where("init_hash IS NULL").Count(&count).Error)
	assert.Equal(int64(3), count)

	assert.NoError(workers.NewHashWorker(db).Work(metadataJobs[jobs.JobMetadataHash]))

	assert.NoError(db.Model(&paths.PathMetadata{}).Where("init_hash IS NOT NULL").Count(&count).Error)
	assert.Equal(int64(3), count)

}
