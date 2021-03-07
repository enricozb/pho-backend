package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Job struct {
	ID        uuid.UUID
	Status    JobStatus
	Kind      JobKind
	CreatedAt time.Time
	UpdatedAt time.Time
	ImportID  uuid.UUID

	Import Import
}

type Import struct {
	ID        uuid.UUID     `gorm:"type:uuid"`
	Opts      ImportOptions `gorm:"-"`
	Status    ImportStatus  `gorm:"default:'NOT_STARTED'"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// OptsJSON should not be set manually, it is used as an intermediate
	// between gorm and unmarshaling to Import.Opts.
	OptsJSON []byte `gorm:"column:opts"`
}

type ImportOptions struct {
	Paths []string `json:"paths"`
}

// BeforeSave creates a new uuid.
func (i *Import) BeforeSave(tx *gorm.DB) error {
	i.ID = uuid.New()
	return nil
}

// BeforeSave creates a new uuid.
func (job *Job) BeforeSave(tx *gorm.DB) error {
	job.ID = uuid.New()
	return nil
}

// BeforeCreate marshals Import.Opts into Import.OptsJSON.
func (i *Import) BeforeCreate(tx *gorm.DB) error {
	if len(i.OptsJSON) != 0 {
		return errors.New("Import.OptsJSON should not be set manually")
	}

	bytes, err := json.Marshal(i.Opts)
	if err != nil {
		return fmt.Errorf("marshall: %v", err)
	}

	i.OptsJSON = bytes
	return nil
}

// AfterFine unmarshals Import.OptsJSON into Import.Opts.
func (i *Import) AfterFind(tx *gorm.DB) error {
	return json.Unmarshal(i.OptsJSON, &i.Opts)
}
