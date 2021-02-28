package files_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
)

func setup(t *testing.T) (*require.Assertions, *sqlx.DB, files.Dao, func()) {
	assert := require.New(t)
	db, cleanup := testutil.MockDB(t)
	dao := files.NewDao(db)

	return assert, db, dao, cleanup
}

func TestFiles_AddDuplicateFiles(t *testing.T) {
	assert, db, dao, cleanup := setup(t)
	defer cleanup()

	importID := testutil.MockImport(t, db)

	numFiles := 1000
	filesToInsert := make([]files.File, numFiles)

	for i := range filesToInsert {
		filesToInsert[i] = files.File{
			ImportID:  importID,
			Kind:      files.ImageKind,
			Timestamp: time.Now(),
			InitHash:  []byte("random hash"),
		}
	}

	fileIDs, err := dao.AddFiles(filesToInsert)
	assert.NoError(err, "add files")

	assert.NotEqual(fileIDs[0], uuid.Nil)
	for _, fileID := range fileIDs[1:] {
		assert.Equal(fileID, uuid.Nil)
	}

	assert.Equal(1, testutil.NumRows(t, db, "files"))
}
