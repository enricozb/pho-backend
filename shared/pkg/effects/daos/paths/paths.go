package paths

import (
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/lib/reflect"
)

type PathID = uuid.UUID

type Path struct {
	ID       PathID        `db:"id"`
	ImportID jobs.ImportID `db:"import_id"`
	Path     string        `db:"path"`
}

type PathMetadata struct {
	PathID    PathID         `db:"path_id"`
	Kind      files.FileKind `db:"kind"`
	Timestamp time.Time      `db:"timestamp"`
	InitHash  []byte         `db:"init_hash"`
	LiveID    []byte         `db:"live_id"`
}

var _ Dao = &dao{}

type Dao interface {
	// Paths returns a slice of all paths for a given import.
	Paths(importID jobs.ImportID) ([]Path, error)

	// AddPaths inserts a slice paths returning the generated PathID's.
	AddPaths(importID jobs.ImportID, paths []string) ([]PathID, error)

	// SetKind sets the kind for a path.
	SetKind(pathID PathID, kind files.FileKind) error

	// SetTimestamp sets the timestamp for a path.
	SetTimestamp(pathID PathID, timestamp time.Time) error

	// SetInitHash sets the initial hash for a path.
	SetInitHash(pathID PathID, hash []byte) error

	// SetLiveID sets the iOS "live" photo id for a path.
	SetLiveID(pathID PathID, liveID []byte) error
}

type dao struct {
	db *sqlx.DB
}

func NewDao(conn *sqlx.DB) *dao {
	return &dao{conn}
}

// Paths returns a slice of all paths for a given import.
func (d *dao) Paths(importID jobs.ImportID) (paths []Path, err error) {
	q, args, err := squirrel.
		Select("*").
		From("paths").
		Where("import_id = ?", importID).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %v", err)
	}

	return paths, d.db.Select(&paths, q, args...)
}

// AddPaths inserts a slice paths returning the generated PathID's.
func (d *dao) AddPaths(importID jobs.ImportID, paths []string) (pathIDs []PathID, err error) {
	for _, path := range paths {
		pathToInsert := Path{
			ID:       uuid.New(),
			ImportID: importID,
			Path:     path,
		}

		q, args, err := squirrel.
			Insert("paths").
			Columns(reflect.Tags(pathToInsert, "db")...).
			Values(reflect.Values(pathToInsert, "db")...).
			ToSql()

		if err != nil {
			return nil, fmt.Errorf("build query: %v", err)
		}

		if _, err := d.db.Exec(q, args...); err != nil {
			return nil, fmt.Errorf("insert path: %v", err)
		}

		pathIDs = append(pathIDs, pathToInsert.ID)
	}

	return pathIDs, nil
}

// SetKind sets the kind for a path.
func (d *dao) SetKind(pathID PathID, kind files.FileKind) error {
	return d.setMetadata(pathID, "kind", kind)
}

// SetTimestamp sets the timestamp for a path.
func (d *dao) SetTimestamp(pathID PathID, timestamp time.Time) error {
	return d.setMetadata(pathID, "timestamp", timestamp)
}

// SetInitHash sets the initial hash for a path.
func (d *dao) SetInitHash(pathID PathID, initHash []byte) error {
	return d.setMetadata(pathID, "init_hash", initHash)
}

// SetLiveID sets the iOS "live" photo id for a path.
func (d *dao) SetLiveID(pathID PathID, liveID []byte) error {
	return d.setMetadata(pathID, "live_id", liveID)
}

func (d *dao) setMetadata(pathID PathID, metadataColumn string, value interface{}) error {
	q, args, err := squirrel.
		Insert("path_metadata").
		Columns("path_id", metadataColumn).
		Values(pathID, value).
		Suffix(fmt.Sprintf("ON CONFLICT(path_id) DO UPDATE SET %s = ?", metadataColumn), value).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query: %v", err)
	}

	_, err = d.db.Exec(q, args...)
	return err
}
