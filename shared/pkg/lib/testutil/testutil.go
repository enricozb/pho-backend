package testutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"gorm.io/gorm"

	"github.com/stretchr/testify/require"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/db"
)

// get media fixtures relative to this file
var _, currentFile, _, _ = runtime.Caller(0)
var MediaFixturesPath = filepath.Join(filepath.Dir(currentFile), ".media_fixtures")

const (
	NumFilesInFixture       int64 = 6
	NumUniqueFilesInFixture int64 = 5
)

func MockDB(t *testing.T) (mockDB *gorm.DB, cleanup func()) {
	assert := require.New(t)

	tmpdir, err := ioutil.TempDir("", "pho-tests-")
	assert.NoError(err, "create tempdir")

	mockDB, err = db.NewDB(tmpdir, "mock-db")
	assert.NoError(err, "new mock db")

	assert.NoError(db.Migrate(mockDB), "migrate mock db")

	return mockDB, func() { os.RemoveAll(tmpdir) }
}

func MockImport(t *testing.T, db *gorm.DB) jobs.Import {
	return MockImportWithOptions(t, db, jobs.ImportOptions{})
}

func MockImportWithOptions(t *testing.T, db *gorm.DB, opts jobs.ImportOptions) jobs.Import {
	assert := require.New(t)
	importEntry := jobs.Import{Opts: opts}
	assert.NoError(db.Create(&importEntry).Error)

	return importEntry

}
