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
	if err := w.db.Find(&importEntry, job.ImportID).Error; err != nil {
		return fmt.Errorf("find import: %v", err)
	}

	if err := importEntry.SetStatus(w.db, jobs.ImportStatusScan); err != nil {
		return fmt.Errorf("set import status: %v", err)
	}

	scannedPaths, err := w.walkPaths(importEntry)
	if err != nil {
		return fmt.Errorf("walk paths: %v", err)
	}

	if err := w.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&scannedPaths).Error; err != nil {
		return fmt.Errorf("add paths: %v", err)
	}

	if _, err := jobs.PushJob(w.db, importEntry.ID, jobs.JobMetadata); err != nil {
		return fmt.Errorf("push job: %v", err)
	}

	return nil
}

// walkPaths returns every supported path under `importEntry.Opts.Paths`.
func (w *scanWorker) walkPaths(importEntry jobs.Import) (supportedPaths []paths.Path, err error) {
	computePath := func(path string) (bool, paths.Path) {
		if isSupported, kind, mimetype := file.Kind(path); isSupported {
			return true, paths.Path{ImportID: importEntry.ID, Path: path, Kind: kind, Mimetype: mimetype}
		}
		return false, paths.Path{}
	}

	for _, path := range importEntry.Opts.Paths {
		if !filepath.IsAbs(path) {
			continue
		}

		if file.IsDir(path) {
			err := godirwalk.Walk(path, &godirwalk.Options{
				Callback: func(path string, de *godirwalk.Dirent) error {
					if !de.IsDir() {
						if shouldAdd, pathEntry := computePath(path); shouldAdd {
							supportedPaths = append(supportedPaths, pathEntry)
						}
					}

					return nil
				},
				Unsorted: true,
			})

			if err != nil {
				return nil, fmt.Errorf("walk '%s': %v", path, err)
			}
		} else if shouldAdd, pathEntry := computePath(path); shouldAdd {
			supportedPaths = append(supportedPaths, pathEntry)
		}
	}

	return supportedPaths, nil
}
