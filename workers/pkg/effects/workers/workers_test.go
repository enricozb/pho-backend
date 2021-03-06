package workers_test

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/enricozb/pho/shared/pkg/effects/daos"
	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
)

func setup(t *testing.T) (*require.Assertions, *sqlx.DB, daos.Dao, func()) {
	assert := require.New(t)
	db, cleanup := testutil.MockDB(t)
	dao := daos.NewDao(db)

	return assert, db, dao, cleanup
}

func assertDidSetImportStatus(t *testing.T, dao jobs.Dao, importID jobs.ImportID, expectedStatus jobs.Status) {
	assert := require.New(t)

	actualStatus, err := dao.GetImportStatus(importID)
	assert.NoError(err, "get import status")

	assert.Equal(expectedStatus, actualStatus)
}
