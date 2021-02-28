package albums_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/enricozb/pho/shared/pkg/effects/daos/albums"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
)

func TestAlbums_NewAlbum(t *testing.T) {
	assert := assert.New(t)

	db, cleanup := testutil.MockDB(t)
	defer cleanup()

	dao := albums.NewDao(db)
	_, err := dao.NewAlbum("test-album", uuid.Nil)
	assert.NoError(err, "new album")

	assert.Equal(1, testutil.NumRows(t, db, "albums"))
}
