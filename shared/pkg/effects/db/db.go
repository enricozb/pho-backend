package db

import (
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB(dir, name string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(filepath.Join(dir, name)), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
}
