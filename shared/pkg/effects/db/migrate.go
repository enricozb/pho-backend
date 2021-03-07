package db

import (
	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
)

var Tables = []interface{}{
	&jobs.Job{},
	&jobs.Import{},
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(Tables...)
}
