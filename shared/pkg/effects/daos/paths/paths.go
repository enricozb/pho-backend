package paths

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
)

type Path struct {
	ID   uuid.UUID
	Path string `gorm:"unique; not null"`

	// exif metadata in json
	EXIFMetadata []byte

	// non-exif metadata
	Kind     files.FileKind
	Mimetype string
	InitHash []byte

	ImportID uuid.UUID
	Import   jobs.Import
}

func (path *Path) BeforeCreate(tx *gorm.DB) error {
	if path.ID == uuid.Nil {
		path.ID = uuid.New()
	}

	if path.Path == "" {
		return errors.New("Path.Path is required")
	}

	return nil
}
