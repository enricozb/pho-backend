package albums_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/enricozb/pho/shared/pkg/effects/daos/albums"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
)

func setup(t *testing.T) (*require.Assertions, *sqlx.DB, albums.Dao, func()) {
	assert := require.New(t)
	db, cleanup := testutil.MockDB(t)
	dao := albums.NewDao(db)

	return assert, db, dao, cleanup
}

func TestAlbums_NewAlbum(t *testing.T) {
	assert, db, dao, cleanup := setup(t)
	defer cleanup()

	assert.Equal(0, testutil.NumRows(t, db, "albums"))

	_, err := dao.NewAlbum("test-album", uuid.Nil)
	assert.NoError(err, "new album")

	assert.Equal(1, testutil.NumRows(t, db, "albums"))
}
