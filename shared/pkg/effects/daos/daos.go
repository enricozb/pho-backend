package daos

import (
	"github.com/jmoiron/sqlx"

	"github.com/enricozb/pho/shared/pkg/effects/daos/albums"
	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
)

type albumsDao albums.Dao
type filesDao files.Dao
type jobsDao jobs.Dao
type pathsDao paths.Dao

type dao struct {
	albumsDao
	filesDao
	jobsDao
	pathsDao
}

var _ Dao = &dao{}

type Dao interface {
	albums.Dao
	files.Dao
	jobs.Dao
	paths.Dao
}

func NewDao(conn *sqlx.DB) *dao {
	return &dao{
		albums.NewDao(conn),
		files.NewDao(conn),
		jobs.NewDao(conn),
		paths.NewDao(conn),
	}
}
