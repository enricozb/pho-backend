package albums

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AlbumID = uuid.UUID

type Album struct {
	ID   AlbumID
	Name string

	ParentAlbumID *AlbumID
	ChildAlbums   []*Album `gorm:"foreignKey:ParentAlbumID"`
}

func (a *Album) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.New()
	return nil
}
