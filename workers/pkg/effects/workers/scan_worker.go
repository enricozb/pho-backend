package workers

import (
	"fmt"
	"path/filepath"

	"github.com/karrick/godirwalk"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
	"github.com/enricozb/pho/shared/pkg/lib/file"
	"github.com/enricozb/pho/workers/pkg/lib/worker"
)

type scanWorker struct {
	db *gorm.DB
}

var _ worker.Worker = &scanWorker{}

func NewScanWorker(db *gorm.DB) *scanWorker {
	return &scanWorker{db: db}
}

func (w *scanWorker) Work(job jobs.Job) error {
	importEntry := jobs.Import{}
	if err := w.db.First(&importEntry, job.ImportID).Error; err != nil {
		return fmt.Errorf("find import: %v", err)
	}

	if err := importEntry.SetStatus(w.db, jobs.ImportStatusScan); err != nil {
		return fmt.Errorf("set import status: %v", err)
	}

	paths, err := w.walkPaths(importEntry)
	if err != nil {
		return fmt.Errorf("walk paths: %v", err)
	}

	if err := w.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&paths).Error; err != nil {
		return fmt.Errorf("add paths: %v", err)
	}

	if _, err := jobs.PushJob(w.db, importEntry.ID, jobs.JobMetadata); err != nil {
		return fmt.Errorf("push job: %v", err)
	}

	return nil
}

// walkPaths returns all file paths under `importEntry.Opts.Paths`.
func (w *scanWorker) walkPaths(importEntry jobs.Import) (supportedPaths []paths.Path, err error) {

	computePath := func(path string) paths.Path {
		isSupported, kind, mimetype := file.Kind(path)
		if isSupported {
			return paths.Path{ImportID: importEntry.ID, Path: path, Kind: kind, Mimetype: mimetype}
		}
		return paths.Path{ImportID: importEntry.ID, Path: path, DiscardReason: fmt.Sprintf("mimetype '%s' is unsupported", mimetype)}
	}

	for _, path := range importEntry.Opts.Paths {
		if !filepath.IsAbs(path) {
			continue
		}

		if file.IsDir(path) {
			err := godirwalk.Walk(path, &godirwalk.Options{
				Callback: func(path string, de *godirwalk.Dirent) error {
					if !de.IsDir() {
						supportedPaths = append(supportedPaths, computePath(path))
					}

					return nil
				},
				Unsorted: true,
			})

			if err != nil {
				return nil, fmt.Errorf("walk '%s': %v", path, err)
			}
		} else {
			supportedPaths = append(supportedPaths, computePath(path))
		}
	}

	return supportedPaths, nil
}
