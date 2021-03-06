package workers

import (
	"fmt"
	"path/filepath"

	"github.com/karrick/godirwalk"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
	"github.com/enricozb/pho/shared/pkg/lib/file"
	"github.com/enricozb/pho/workers/pkg/lib/worker"
)

type scanWorkerDao interface {
	jobs.Dao
	paths.Dao
}

type scanWorker struct {
	dao scanWorkerDao
}

var _ worker.Worker = &scanWorker{}

func NewScanWorker(dao scanWorkerDao) *scanWorker {
	return &scanWorker{dao: dao}
}

func (w *scanWorker) Work(jobID jobs.JobID) error {
	importID, opts, err := getJobInfo(w.dao, jobID)
	if err != nil {
		return fmt.Errorf("get job info: %v", err)
	}

	paths, err := w.walkPaths(opts.Paths)
	if err != nil {
		return fmt.Errorf("walk paths: %v", err)
	}

	if _, err := w.dao.AddPaths(importID, paths); err != nil {
		return fmt.Errorf("add paths: %v", err)
	}

	return nil
}

// walkPaths returns every supporeted path under `inPaths`.
func (w *scanWorker) walkPaths(inPaths []string) (paths []string, err error) {
	fmt.Printf("trying to find paths in %s\n", paths)

	for _, path := range inPaths {
		if !filepath.IsAbs(path) {
			continue
		}

		if file.IsSupported(path) {
			paths = append(paths, path)
			continue
		}

		if file.IsDir(path) {
			err := godirwalk.Walk(path, &godirwalk.Options{
				Callback: func(path string, de *godirwalk.Dirent) error {
					if file.IsSupported(path) {
						paths = append(paths, path)
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

	return paths, nil
}
