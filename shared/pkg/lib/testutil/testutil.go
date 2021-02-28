package testutil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Masterminds/squirrel"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/db"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func MockDB(t *testing.T) (mockDB *sqlx.DB, cleanup func()) {
	assert := assert.New(t)

	tmpdir, err := ioutil.TempDir("", "pho-tests-")
	assert.NoError(err, "create tempdir")

	mockDB, err = db.NewDB(tmpdir, "mock-db")
	assert.NoError(err, "new mock db")

	assert.NoError(db.Migrate(mockDB), "migrate mock db")

	return mockDB, func() { os.RemoveAll(tmpdir) }
}

func MockImport(t *testing.T, db *sqlx.DB) jobs.ImportID {
	importID, err := jobs.NewDao(db).NewImport(jobs.ImportOptions{})
	assert.NoError(t, err, "new import")

	return importID
}

func NumRows(t *testing.T, db *sqlx.DB, tableName string) (count int) {
	assert := assert.New(t)

	q, args, err := squirrel.
		Select("count(*)").
		From(tableName).
		ToSql()

	assert.NoError(err, "build count query")

	assert.NoError(db.Get(&count, q, args...), "query count")
	return count
}
