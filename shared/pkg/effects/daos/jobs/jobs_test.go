package jobs_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
)

func setup(t *testing.T) (*require.Assertions, *gorm.DB, func()) {
	assert := require.New(t)
	db, cleanup := testutil.MockDB(t)

	return assert, db, cleanup
}

func TestJobs_JobSmokeTest(t *testing.T) {
	assert, db, cleanup := setup(t)
	defer cleanup()

	insertedImport := jobs.Import{}
	db.Create(&insertedImport)

	insertedJob := jobs.Job{ImportID: insertedImport.ID, Kind: jobs.JobScan}
	db.Create(&insertedJob)

	foundJob := jobs.Job{ImportID: insertedImport.ID}
	assert.NoError(db.First(&foundJob).Error)

	assert.Equal(jobs.JobStatusNotStarted, insertedJob.Status)
	assert.Equal(insertedJob, foundJob)

	jobs := []jobs.Job{}
	assert.NoError(db.Where("import_id = ?", insertedImport.ID.String()).Find(&jobs).Error)

	assert.Len(jobs, 1)
	assert.Equal(insertedJob, jobs[0])
}

func TestJobs_ImportSmokeTest(t *testing.T) {
	assert, db, cleanup := setup(t)
	defer cleanup()

	insertedImport := jobs.Import{Opts: jobs.ImportOptions{Paths: []string{"a", "b", "c"}}}
	db.Create(&insertedImport)

	foundImport := jobs.Import{ID: uuid.New()}
	assert.EqualError(db.First(&foundImport).Error, "record not found")

	foundImport = jobs.Import{ID: insertedImport.ID}
	assert.NoError(db.First(&foundImport).Error)

	assert.Equal(insertedImport, foundImport)
}
