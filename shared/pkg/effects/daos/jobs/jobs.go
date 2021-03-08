package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ImportID = uuid.UUID
type JobID = uuid.UUID

type Job struct {
	ID        JobID
	Status    JobStatus `gorm:"default:'NOT_STARTED'"`
	Kind      JobKind
	CreatedAt time.Time
	UpdatedAt time.Time
	ImportID  ImportID

	Import Import
}

type Import struct {
	ID        ImportID      `gorm:"type:uuid"`
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

type ImportFailure struct {
	ID       uint
	ImportID ImportID
	Message  string

	Import Import
}

func (job *Job) BeforeCreate(tx *gorm.DB) error {
	job.ID = uuid.New()

	if job.Kind == "" {
		return errors.New("Job.Kind is required")
	}

	return nil
}

func (i *Import) BeforeCreate(tx *gorm.DB) error {
	i.ID = uuid.New()

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

// AfterFind unmarshals Import.OptsJSON into Import.Opts.
func (i *Import) AfterFind(tx *gorm.DB) error {
	return json.Unmarshal(i.OptsJSON, &i.Opts)
}
