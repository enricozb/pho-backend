package db

import (
	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/albums"
	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
)

var Tables = []interface{}{
	&jobs.Job{},
	&jobs.Import{},
	&jobs.ImportFailure{},

	&files.File{},
	&paths.Path{},
	&paths.PathMetadata{},

	&albums.Album{},
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(Tables...)
}
