package files

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
)

type FileID = uuid.UUID

type File struct {
	ID        FileID
	ImportID  jobs.ImportID
	Kind      FileKind
	Timestamp time.Time
	InitHash  []byte `gorm:"unique"`
	ConvHash  []byte `gorm:"unique"`
	LiveID    []byte

	Import jobs.Import
}

func (file *File) BeforeCreate(tx *gorm.DB) error {
	file.ID = uuid.New()

	if file.Kind == "" {
		return errors.New("File.Kind is required")
	}

	return nil
}
