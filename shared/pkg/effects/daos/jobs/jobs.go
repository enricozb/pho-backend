package jobs

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobID = uuid.UUID

type Job struct {
	ID        JobID
	Status    JobStatus `gorm:"default:'NOT_STARTED'"`
	Kind      JobKind   `gorm:"not null"`
	Args      []byte    `gorm:"default:'{}'"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	ImportID  ImportID  `gorm:"not null"`

	Import Import
}

func (job *Job) BeforeCreate(tx *gorm.DB) error {
	if job.ID == uuid.Nil {
		job.ID = uuid.New()
	}

	if job.Kind == "" {
		return errors.New("Job.Kind is required")
	}

	return nil
}

func (job *Job) GetArgs(i interface{}) error {
	return json.Unmarshal(job.Args, i)
}

func (job *Job) SetStatus(db *gorm.DB, status JobStatus) error {
	return db.Model(&Job{}).Where("id = ?", job.ID).Update("status", status).Error
}

func GetJobStatus(db *gorm.DB, jobID JobID) (status JobStatus, err error) {
	var job Job
	return job.Status, db.Where("id = ?", jobID).First(&job).Error
}
