package db

import (
	"fmt"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/config"
)

const defaultDBFileName = "pho.db"

func NewDB(dir, name string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(filepath.Join(dir, name)), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("open: %v", err)
	}

	if err := Migrate(db); err != nil {
		return nil, fmt.Errorf("migrate: %v", err)
	}

	return db, nil
}

func MustDB() *gorm.DB {
	if db, err := NewDB(config.Config.DBDir, defaultDBFileName); err != nil {
		panic(fmt.Errorf("must db: %v", err))
	} else {
		return db
	}
}
