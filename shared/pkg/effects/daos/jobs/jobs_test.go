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
