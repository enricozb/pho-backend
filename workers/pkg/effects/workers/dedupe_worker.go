package workers

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

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
	if err := w.db.First(&importEntry, job.ImportID).Error; err != nil {
		return fmt.Errorf("find import: %v", err)
	}

	if err := importEntry.SetStatus(w.db, jobs.ImportStatusDedupe); err != nil {
		return fmt.Errorf("set import status: %v", err)
	}

	pathsToImport, err := paths.PathsInPipeline(w.db, importEntry.ID)
	if err != nil {
		return fmt.Errorf("get paths: %v", err)
	}

	filesToImport := make([]files.File, len(pathsToImport))
	for i, path := range pathsToImport {
		filesToImport[i].ID = path.ID
		filesToImport[i].ImportID = path.ImportID
		filesToImport[i].Kind = path.Kind
		filesToImport[i].Timestamp = path.EXIFMetadata.Timestamp
		filesToImport[i].LiveID = path.EXIFMetadata.LiveID
		filesToImport[i].InitHash = path.InitHash
		filesToImport[i].Width = path.EXIFMetadata.Width
		filesToImport[i].Height = path.EXIFMetadata.Height
	}

	if err := w.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&filesToImport).Error; err != nil {
		return fmt.Errorf("insert files: %v", err)
	}

	if _, err := jobs.PushJob(w.db, importEntry.ID, jobs.JobConvert); err != nil {
		return fmt.Errorf("push job: %v", err)
	}

	return nil
}
