package jobs

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

func (job *Job) BeforeCreate(tx *gorm.DB) error {
	job.ID = uuid.New()

	if job.Kind == "" {
		return errors.New("Job.Kind is required")
	}

	return nil
}
