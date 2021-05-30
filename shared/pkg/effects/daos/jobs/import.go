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

func StartImport(db *gorm.DB, importOptions ImportOptions) error {
	importEntry := Import{Opts: importOptions}

	if err := db.Create(&importEntry).Error; err != nil {
		return fmt.Errorf("create import: %v", err)
	}

	if _, err := PushJob(db, importEntry.ID, JobScan); err != nil {
		return fmt.Errorf("push job: %v", err)
	}

	return nil
}

type ImportFailure struct {
	ID       uint
	ImportID ImportID
	Message  string

	Import Import
}

func (i *Import) BeforeSave(tx *gorm.DB) (err error) {
	if len(i.OptsJSON) != 0 {
		return errors.New("Import.OptsJSON should not be set manually")
	}

	if i.OptsJSON, err = json.Marshal(i.Opts); err != nil {
		return fmt.Errorf("marshal: %v", err)
	}

	return nil
}

func (i *Import) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}

	return nil
}

// AfterFind unmarshals Import.OptsJSON into Import.Opts.
func (i *Import) AfterFind(tx *gorm.DB) error {
	// clear OptsJSON so it cannot be read after a find, it will be re-marshaled on save
	optsJSON := i.OptsJSON
	i.OptsJSON = []byte("")

	return json.Unmarshal(optsJSON, &i.Opts)
}

func (i *Import) AfterSave(tx *gorm.DB) error {
	// clear OptsJSON so it cannot be read after saving
	i.OptsJSON = []byte("")

	return nil
}

func (i *Import) SetStatus(db *gorm.DB, status ImportStatus) error {
	return db.Model(&Import{}).Where("id = ?", i.ID).Update("status", status).Error
}
