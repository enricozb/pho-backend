package paths

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
)

type Path struct {
	ID   uuid.UUID
	Path string `gorm:"unique; not null"`

	// auto-unmarshall into EXIFMetadata from EXIFMetadataJSON when fetching
	EXIFMetadata     EXIFMetadata `gorm:"-"`
	EXIFMetadataJSON []byte       `gorm:"column:exif_metadata"`

	// non-exif metadata
	Kind     files.FileKind
	Mimetype string
	InitHash []byte

	// DiscardReason explains why this file was ignored. If not empty, then this path was not imported and no metadata for this path is guaranteed to be set.
	DiscardReason string `gorm:"default:''"`

	ImportID uuid.UUID
	Import   jobs.Import
}

func PathsInPipeline(db *gorm.DB, importID jobs.ImportID) (validPaths []Path, err error) {
	return validPaths, db.Where("import_id = ? AND LENGTH(discard_reason) = 0", importID).Find(&validPaths).Error
}

func (p *Path) BeforeSave(tx *gorm.DB) (err error) {
	if len(p.EXIFMetadataJSON) != 0 {
		return errors.New("Path.EXIFMetadataJSON should not be set manually")
	}

	if p.EXIFMetadata.Path == "" {
		p.EXIFMetadataJSON = []byte("{}")
	} else if p.EXIFMetadataJSON, err = json.Marshal(p.EXIFMetadata); err != nil {
		return fmt.Errorf("marshal: %v", err)
	}

	return nil
}

func (p *Path) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	return nil
}

// AfterFind unmarshals Path.EXIFMetadataJSON into Path.EXIFMetadata.
func (p *Path) AfterFind(tx *gorm.DB) error {
	// clear EXIFMetadataJSON so it cannot be read after a find, it will be re-marshaled on save
	exifMetadataJSON := p.EXIFMetadataJSON
	p.EXIFMetadataJSON = []byte("")

	return json.Unmarshal(exifMetadataJSON, &p.EXIFMetadata)
}

func (p *Path) AfterSave(tx *gorm.DB) (err error) {
	// clear EXIFMetadataJSON so it cannot be read after saving
	p.EXIFMetadataJSON = []byte("")

	return nil
}
