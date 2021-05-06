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

type Import struct {
	ID        ImportID      `gorm:"type:uuid"`
	Opts      ImportOptions `gorm:"-"`
	Status    ImportStatus  `gorm:"default:'NOT_STARTED'"`
	CreatedAt time.Time     `gorm:"not null"`
	UpdatedAt time.Time     `gorm:"not null"`

	// OptsJSON should not be set manually, it is used as an intermediate between gorm and unmarshaling to Import.Opts.
	OptsJSON []byte `gorm:"column:opts;not null"`
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

func (i *Import) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}

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

func (i *Import) SetStatus(db *gorm.DB, status ImportStatus) error {
	return db.Model(&Import{}).Where("id = ?", i.ID).Update("status", status).Error
}
