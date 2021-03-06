package files

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/exif"
	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
)

type FileID = uuid.UUID

type File struct {
	ID       FileID
	ImportID jobs.ImportID `gorm:"not null"`
	Kind     FileKind      `gorm:"not null"`
	InitHash []byte        `gorm:"unique;not null"`

	exif.EXIFMetadata

	Extension string `gorm:"default:NULL"`

	Import jobs.Import
}

func (file *File) BeforeCreate(tx *gorm.DB) error {
	if file.ID == uuid.Nil {
		file.ID = uuid.New()
	}

	if file.Kind == "" {
		return errors.New("File.Kind is required")
	}

	return nil
}
