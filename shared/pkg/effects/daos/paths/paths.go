package paths

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
)

type Path struct {
	ID       uuid.UUID
	ImportID uuid.UUID
	Path     string

	Import jobs.Import
}

type PathMetadata struct {
	PathID    uuid.UUID
	Kind      files.FileKind
	Timestamp sql.NullTime
	InitHash  []byte
	LiveID    []byte

	Path Path
}

func (path *Path) BeforeCreate(tx *gorm.DB) error {
	path.ID = uuid.New()

	if path.Path == "" {
		return errors.New("Path.Path is required")
	}

	return nil
}
