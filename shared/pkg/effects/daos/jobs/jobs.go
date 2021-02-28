package jobs

import (
	"encoding/json"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ImportID = uuid.UUID
type JobID = uuid.UUID

type Job struct {
	ID        uuid.UUID `db:"id"`
	ImportID  ImportID  `db:"import_id"`
	Kind      JobKind   `db:"kind"`
	CreatedAt time.Time `db:"created_at"`
}

type ImportOptions struct {
	Paths []string `json:"paths"`
}

var _ Dao = &dao{}

type Dao interface {
	// NewImport starts a new import.
	NewImport(opts ImportOptions) (ImportID, error)

	// GetImportStatus retrieves the status for an import.
	GetImportStatus(importID ImportID) (Status, error)

	// SetImportStatus sets the status for an existing import.
	SetImportStatus(importID ImportID, status Status) error

	// PushJob adds a new job to the queue.
	PushJob(importID ImportID, kind JobKind) (JobID, error)

	// PeekJob retrieves a job from the queue, but does not delete it.
	PeekJob(importID ImportID) (Job, error)

	// PopJob retrieves a job from the queue, and deletes it.
	PopJob(importID ImportID) (Job, error)

	// DeleteJob deletes a job from the queue.
	DeleteJob(jobID uuid.UUID) error
}

type dao struct {
	db *sqlx.DB
}

func NewDao(conn *sqlx.DB) *dao {
	return &dao{conn}
}

func (d *dao) NewImport(opts ImportOptions) (importID ImportID, err error) {
	optsBytes, err := json.Marshal(opts)
	if err != nil {
		return uuid.Nil, fmt.Errorf("json marshal: %v", err)
	}

	importID = uuid.New()

	q, args, err := sq.
		Insert("imports").
		Columns("id", "opts").
		Values(importID, optsBytes).
		ToSql()

	if err != nil {
		return uuid.Nil, fmt.Errorf("build query: %v", err)
	}

	_, err = d.db.Exec(q, args...)
	return importID, err
}

// GetImportStatus retrieves the status for an import.
func (d *dao) GetImportStatus(importID ImportID) (status Status, err error) {
	q, args, err := sq.
		Select("status").
		From("imports").
		Where("id = ?", importID).
		ToSql()

	if err != nil {
		return "", fmt.Errorf("build query: %v", err)
	}

	return status, d.db.Get(&status, q, args...)
}

// SetImportStatus sets the status for an existing import.
func (d *dao) SetImportStatus(importID ImportID, status Status) error {
	q, args, err := sq.
		Update("imports").
		Set("status", status).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query: %v", err)
	}

	_, err = d.db.Exec(q, args...)
	return err
}

// PushJob adds a new job to the queue.
func (d *dao) PushJob(importID ImportID, kind JobKind) (jobID JobID, err error) {
	jobID = uuid.New()

	q, args, err := sq.
		Insert("jobs").
		Columns("id", "import_id", "kind").
		Values(jobID, importID, kind).
		ToSql()

	if err != nil {
		return uuid.Nil, fmt.Errorf("build query: %v", err)
	}

	_, err = d.db.Exec(q, args...)
	return jobID, err
}

// PeekJob retrieves a job from the queue, but does not delete it.
func (d *dao) PeekJob(importID ImportID) (job Job, err error) {
	q, args, err := sq.
		Select("*").
		From("jobs").
		Limit(1).
		ToSql()

	if err != nil {
		return Job{}, fmt.Errorf("build query: %v", err)
	}

	return job, d.db.Get(&job, q, args...)
}

// PopJob retrieves a job from the queue, and deletes it.
func (d *dao) PopJob(importID ImportID) (Job, error) {
	job, err := d.PeekJob(importID)
	if err != nil {
		return Job{}, fmt.Errorf("pop job: %w", err)
	}

	return job, d.DeleteJob(job.ID)
}

// DeleteJob deletes a job from the queue.
func (d *dao) DeleteJob(jobID uuid.UUID) error {
	q, args, err := sq.
		Delete("jobs").
		Where("id = ?", jobID).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query: %v", err)
	}

	_, err = d.db.Exec(q, args...)
	return err
}
