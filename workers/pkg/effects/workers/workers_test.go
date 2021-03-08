package workers_test

import (
	"testing"

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

func assertDidSetImportStatus(assert *require.Assertions, db *gorm.DB, importID jobs.ImportID, expectedStatus jobs.ImportStatus) {
	importEntry := jobs.Import{ID: importID}
	assert.NoError(db.Find(&importEntry).Error)
	assert.Equal(expectedStatus, importEntry.Status)
}

func assertDidEnqueueJob(assert *require.Assertions, db *gorm.DB, importID jobs.ImportID, kind jobs.JobKind) {
	assert.NoError(db.Where("import_id = ? AND kind = ?", importID, kind).Find(&jobs.Job{}).Error)
}
