package albums

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
	"github.com/enricozb/pho/shared/pkg/lib/reflect"
)

type AlbumID = uuid.UUID

var _ Dao = &dao{}

type Dao interface {
	// NewAlbum creates a new album under the provided parentID, if any.
	NewAlbum(name string, parentID AlbumID) (AlbumID, error)

	// GetAlbumFiles gets all files for an album.
	GetAlbumFiles(albumID AlbumID) ([]files.File, error)
}

type dao struct {
	db *sqlx.DB
}

func NewDao(conn *sqlx.DB) *dao {
	return &dao{conn}
}

// NewAlbum creates a new album under the provided parentID, if any.
func (d *dao) NewAlbum(name string, parentID AlbumID) (albumID AlbumID, err error) {
	// insert album and get new AlbumID
	q, args, err := sq.
		Insert("albums").
		Columns("name").
		Values(name).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return uuid.Nil, fmt.Errorf("build insert album query: %v", err)
	}

	if err := d.db.Get(&albumID, q, args...); err != nil {
		return uuid.Nil, fmt.Errorf("insert album: %w", err)
	}

	// insert album under parentID
	if parentID == uuid.Nil {
		return albumID, nil
	}

	q, args, err = sq.
		Insert("album_albums").
		Columns("album_id", "child_album_id").
		Values(parentID, albumID).
		ToSql()

	if err != nil {
		return uuid.Nil, fmt.Errorf("build insert album_albums query: %v", err)
	}

	if _, err := d.db.Exec(q, args...); err != nil {
		return uuid.Nil, fmt.Errorf("insert album parent: %w", err)
	}

	return albumID, nil
}

// GetAlbumFiles creates a new album under the provided parentID, if any.
func (d *dao) GetAlbumFiles(albumID AlbumID) (fileIDs []files.File, err error) {
	q, args, err := sq.
		Select(reflect.Tags(files.File{}, "db")...).
		From("album_files").
		Where("album_id = ?", albumID).
		Join("files ON album_files.file_id = files.id").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %v", err)
	}

	return fileIDs, d.db.Select(&fileIDs, q, args...)
}
