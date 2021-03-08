package albums

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
)

type AlbumID = uuid.UUID

type Album struct {
	ID   AlbumID
	Name string

	ParentAlbumID *AlbumID
	ChildAlbums   []*Album     `gorm:"foreignKey:ParentAlbumID"`
	Files         []files.File `gorm:"many2many:album_files"`
}

func (a *Album) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.New()
	return nil
}
