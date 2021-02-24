package files

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/lib/reflect"
)

type FileID = uuid.UUID

type File struct {
	id        FileID        `db:"id"`
	importId  jobs.ImportID `db:"import_id"`
	kind      FileKind      `db:"kind"`
	timestamp time.Time     `db:"timestamp"`
	initHash  []byte        `db:"init_hash"`
	convHash  []byte        `db:"conv_hash"`
	live      uuid.UUID     `db:"live_id"`
}

var _ Dao = &dao{}

type Dao interface {
	// Files lists all files.
	Files() ([]File, error)

	// AddFiles inserts a slice files returning the generated FileID's.
	AddFiles(files []File) ([]FileID, error)
}

type dao struct {
	db *sqlx.DB
}

func NewDao(conn *sqlx.DB) *dao {
	return &dao{conn}
}

// Files lists all files.
func (d *dao) Files() (files []File, err error) {
	q, args, err := sq.
		Select("files").
		Columns("*").
		ToSql()

	if err != nil {
		return files, fmt.Errorf("build query: %v", err)
	}

	return files, d.db.Select(&files, q, args...)
}

// AddFiles inserts a slice files returning the generated FileID's.
func (d *dao) AddFiles(files []File) (fileIDs []FileID, err error) {
	for _, file := range files {
		q, args, err := sq.
			Insert("files").
			Columns(reflect.Tags(file, "db")...).
			Values(reflect.Values(file, "db")...).
			Suffix("RETURNING id").
			ToSql()

		if err != nil {
			return nil, fmt.Errorf("build query: %v", err)
		}

		// TODO(enricozb): check for uniqueness constraint error
		var fileID FileID
		if err := d.db.Get(&fileID, q, args...); err != nil {
			return nil, fmt.Errorf("insert file: %w", err)
		}

		fileIDs = append(fileIDs, fileID)
	}

	return fileIDs, nil
}
