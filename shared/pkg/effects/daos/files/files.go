package files

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/lib/reflect"
)

type FileID = uuid.UUID

type File struct {
	ID        FileID        `db:"id"`
	ImportID  jobs.ImportID `db:"import_id"`
	Kind      FileKind      `db:"kind"`
	Timestamp time.Time     `db:"timestamp"`
	InitHash  []byte        `db:"init_hash"`
	ConvHash  []byte        `db:"conv_hash"`
	LiveID    []byte        `db:"live_id"`
}

var _ Dao = &dao{}

type Dao interface {
	// Files lists all files.
	Files() ([]File, error)

	// AddFiles inserts a slice files returning the generated FileID's, where some FileIDs are uuid.Nil if duplicates were found.
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
		Select("*").
		From("files").
		ToSql()

	if err != nil {
		return files, fmt.Errorf("build query: %v", err)
	}

	return files, d.db.Select(&files, q, args...)
}

// AddFiles inserts a slice files returning the generated FileID's.
func (d *dao) AddFiles(files []File) (fileIDs []FileID, err error) {
	for _, file := range files {
		file.ID = uuid.New()

		q, args, err := sq.
			Insert("files").
			Columns(reflect.Tags(file, "db")...).
			Values(reflect.Values(file, "db")...).
			ToSql()

		if err != nil {
			return nil, fmt.Errorf("build query: %v", err)
		}

		_, err = d.db.Exec(q, args...)

		// if an unique constraint violation occurs during insert, don't set a uuid for that file.
		if sqle, ok := err.(sqlite3.Error); ok && sqle.ExtendedCode == sqlite3.ErrConstraintUnique {
			file.ID = uuid.Nil
		} else if err != nil {
			return nil, fmt.Errorf("insert file: %v", err)
		}

		fileIDs = append(fileIDs, file.ID)
	}

	return fileIDs, nil
}
