package workers

import (
	"fmt"
	"path/filepath"

	"github.com/karrick/godirwalk"
	"gorm.io/gorm"

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

	if err := w.db.Create(&scannedPaths).Error; err != nil {
		return fmt.Errorf("add paths: %v", err)
	}

	pathMetadatas := make([]paths.PathMetadata, len(scannedPaths))
	for i := range pathMetadatas {
		pathMetadatas[i].PathID = scannedPaths[i].ID
	}

	if err := w.db.Create(&pathMetadatas).Error; err != nil {
		return fmt.Errorf("add paths: %v", err)
	}

	if _, err := jobs.PushJob(w.db, importEntry.ID, jobs.JobMetadata); err != nil {
		return fmt.Errorf("push job: %v", err)
	}

	return nil
}

// walkPaths returns every supported path under `importEntry.Opts.Paths`.
func (w *scanWorker) walkPaths(importEntry jobs.Import) (supportedPaths []paths.Path, err error) {
	for _, path := range importEntry.Opts.Paths {
		if !filepath.IsAbs(path) {
			continue
		}

		if file.IsSupported(path) {
			supportedPaths = append(supportedPaths, paths.Path{ImportID: importEntry.ID, Path: path})
			continue
		}

		if file.IsDir(path) {
			err := godirwalk.Walk(path, &godirwalk.Options{
				Callback: func(path string, de *godirwalk.Dirent) error {
					if file.IsSupported(path) {
						supportedPaths = append(supportedPaths, paths.Path{ImportID: importEntry.ID, Path: path})
					}

					return nil
				},
				Unsorted: true,
			})

			if err != nil {
				return nil, fmt.Errorf("walk '%s': %v", path, err)
			}
		}
	}

	return supportedPaths, nil
}
