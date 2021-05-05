package workers

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
	"github.com/enricozb/pho/workers/pkg/lib/worker"
)

type dedupeWorker struct {
	db *gorm.DB
}

var _ worker.Worker = &dedupeWorker{}

func NewDedupeWorker(db *gorm.DB) *dedupeWorker {
	return &dedupeWorker{db: db}
}

func (w *dedupeWorker) Work(job jobs.Job) error {
	importEntry := jobs.Import{}
	if err := w.db.Find(&importEntry, job.ImportID).Error; err != nil {
		return fmt.Errorf("find import: %v", err)
	}

	var pathsToImport []paths.Path
	if err := w.db.Where("import_id = ?", importEntry.ID).Find(&pathsToImport).Error; err != nil {
		return fmt.Errorf("get paths: %v", err)
	}

	filesToImport := make([]files.File, len(pathsToImport))
	for i, path := range pathsToImport {
		if !path.Timestamp.Valid {
			return fmt.Errorf("path without timestamp")
		}

		filesToImport[i].ID = path.ID
		filesToImport[i].Kind = path.Kind
		filesToImport[i].Timestamp = path.Timestamp.Time
		filesToImport[i].InitHash = path.InitHash
		filesToImport[i].LiveID = path.LiveID
	}

	if _, err := jobs.PushJob(w.db, importEntry.ID, jobs.JobConvert); err != nil {
		return fmt.Errorf("push job: %v", err)
	}

	return nil
}
