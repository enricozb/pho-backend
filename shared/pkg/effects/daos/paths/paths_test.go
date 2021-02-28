package paths_test

import (
	"crypto/md5"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
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
	insertedPaths := seedRandomPaths(t, dao, importID, numPaths)

	assert.Equal(numPaths, testutil.NumRows(t, db, "paths"))
	assert.Equal(0, testutil.NumRows(t, db, "path_metadata"))

	// check that inserted paths match the expected path strings
	retrivedPaths, err := dao.Paths(importID)
	assert.NoError(err, "paths")

	assert.ElementsMatch(insertedPaths, retrivedPaths)
}

func TestPaths_PathMetadata(t *testing.T) {
	assert, db, dao, cleanup := setup(t)
	defer cleanup()

	importID := testutil.MockImport(t, db)

	numPaths := 100
	insertedPaths := seedRandomPaths(t, dao, importID, numPaths)

	pathMetadatas := make([]paths.PathMetadata, numPaths)
	for i, path := range insertedPaths {
		pathMetadatas[i] = randomPathMetadata(path.ID)

		assert.NoError(dao.SetKind(path.ID, pathMetadatas[i].Kind), "set kind")
		assert.NoError(dao.SetTimestamp(path.ID, pathMetadatas[i].Timestamp), "set timestamp")
		assert.NoError(dao.SetInitHash(path.ID, pathMetadatas[i].InitHash), "set init hash")
		assert.NoError(dao.SetLiveID(path.ID, pathMetadatas[i].LiveID), "set live id")
	}

	assert.Equal(numPaths, testutil.NumRows(t, db, "paths"))
	assert.Equal(numPaths, testutil.NumRows(t, db, "path_metadata"))
}

func seedRandomPaths(t *testing.T, dao paths.Dao, importID jobs.ImportID, numPaths int) []paths.Path {
	pathStructs := make([]paths.Path, numPaths)
	pathStrings := make([]string, numPaths)
	for i := range pathStrings {
		// using uuids as path names for uniqueness, not a normal use-case
		pathStrings[i] = uuid.New().String()

		pathStructs[i] = paths.Path{
			ImportID: importID,
			Path:     pathStrings[i],
		}
	}

	pathIDs, err := dao.AddPaths(importID, pathStrings)
	assert.NoError(t, err, "add paths")

	for i, pathID := range pathIDs {
		pathStructs[i].ID = pathID
	}

	return pathStructs
}

// randomPathMetadata creates a random PathMetadata.
func randomPathMetadata(pathID paths.PathID) paths.PathMetadata {
	kind := files.ImageKind
	if rand.Float32() > 0.5 {
		kind = files.VideoKind
	}

	initHash := md5.Sum([]byte(uuid.New().String()))

	return paths.PathMetadata{
		PathID:    pathID,
		Kind:      kind,
		Timestamp: time.Now(),
		InitHash:  initHash[:],
		LiveID:    []byte(uuid.New().String()),
	}
}
