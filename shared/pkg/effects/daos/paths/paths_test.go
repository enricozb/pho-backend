package paths_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
)

func setup(t *testing.T) (*assert.Assertions, *sqlx.DB, paths.Dao, func()) {
	assert := assert.New(t)
	db, cleanup := testutil.MockDB(t)
	dao := paths.NewDao(db)

	return assert, db, dao, cleanup
}

func TestPaths_AddPaths(t *testing.T) {
	assert, db, dao, cleanup := setup(t)
	defer cleanup()

	importID := testutil.MockImport(t, db)

	assert.Equal(0, testutil.NumRows(t, db, "paths"))

	numPaths := 100
	paths := make([]string, numPaths)
	for i := range paths {
		// using uuids as path names for uniqueness, not a normal use-case
		paths[i] = uuid.New().String()
	}

	_, err := dao.AddPaths(importID, paths)
	assert.NoError(err, "add paths")

	assert.Equal(numPaths, testutil.NumRows(t, db, "paths"))

	// check that inserted paths match the expected path strings
	insertedPaths, err := dao.Paths(importID)
	insertedPathStrings := make([]string, len(insertedPaths))
	for i, path := range insertedPaths {
		insertedPathStrings[i] = path.Path
	}

	assert.ElementsMatch(paths, insertedPathStrings)
}
