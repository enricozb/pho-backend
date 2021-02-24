package jobs

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ImportID = uuid.UUID

type Job struct {
	id        uuid.UUID `db:"id"`
	import_id ImportID  `db:"import_id"`
	kind      JobKind   `db:"kind"`
}

var _ Dao = &dao{}

type Dao interface {
	// NewImport starts a new import.
	NewImport(dirs []string) (ImportID, error)

	// GetImportStatus retrieves the status for an import.
	GetImportStatus(importID ImportID) (Status, error)

	// SetImportStatus sets the status for an existing import.
	SetImportStatus(importID ImportID, status Status) error

	// PushJob adds a new job to the queue.
	PushJob(importID ImportID, kind JobKind) error

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

func (d *dao) NewImport(dirs []string) (importID ImportID, err error) {
	q, args, err := sq.
		Insert("imports").
		Columns("dirs").
		Values(dirs).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return uuid.Nil, fmt.Errorf("build query: %v", err)
	}

	return importID, d.db.Get(&importID, q, args...)
}

// GetImportStatus retrieves the status for an import.
func (d *dao) GetImportStatus(importID ImportID) (status Status, err error) {
	q, args, err := sq.
		Select("imports").
		Columns("status").
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
func (d *dao) PushJob(importID ImportID, kind JobKind) error {
	q, args, err := sq.
		Insert("jobs").
		Columns("import_id", "kind").
		Values(importID, kind).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query: %v", err)
	}

	_, err = d.db.Exec(q, args...)
	return err
}

// PeekJob retrieves a job from the queue, but does not delete it.
func (d *dao) PeekJob(importID ImportID) (job Job, err error) {
	q, args, err := sq.
		Select("jobs").
		Columns("*").
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

	return job, d.DeleteJob(job.id)
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
