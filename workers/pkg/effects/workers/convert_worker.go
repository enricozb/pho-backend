package workers

import (
	"fmt"
	"path/filepath"

	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/config"
	"github.com/enricozb/pho/shared/pkg/effects/converter"
	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
	"github.com/enricozb/pho/workers/pkg/lib/worker"
)

type ConvertWorker struct {
	db *gorm.DB
}

var _ worker.Worker = &ConvertWorker{}

func NewConvertWorker(db *gorm.DB) *ConvertWorker {
	return &ConvertWorker{db: db}
}

func (w *ConvertWorker) Work(job jobs.Job) error {
	importEntry := jobs.Import{}
	if err := w.db.First(&importEntry, job.ImportID).Error; err != nil {
		return fmt.Errorf("find import: %v", err)
	}

	if err := importEntry.SetStatus(w.db, jobs.ImportStatusConvert); err != nil {
		return fmt.Errorf("set import status: %v", err)
	}

	var filesToImport []files.File
	if err := w.db.Where("import_id = ? AND kind = ?", importEntry.ID, files.ImageKind).Find(&filesToImport).Error; err != nil {
		return fmt.Errorf("get image files: %v", err)
	}

	// convert the files
	converter := converter.NewMediaConverter()

	for _, file := range filesToImport {
		var path paths.Path
		if err := w.db.Where("id = ?", file.ID).First(&path).Error; err != nil {
			return fmt.Errorf("get path: %v", err)
		}

		src := path.Path
		dst := destPath(file)
		if err := converter.Convert(src, dst, path.Mimetype); err != nil {
			return fmt.Errorf("convert: %v", err)
		}
	}

	if err := converter.Finish(); err != nil {
		return fmt.Errorf("finish: %v", err)
	}

	if _, err := jobs.PushJob(w.db, importEntry.ID, jobs.JobCleanup); err != nil {
		return fmt.Errorf("push job: %v", err)
	}

	return nil
}

// destPath returns the destination path after conversion for a file, _without the new extension_.
func destPath(file files.File) string {
	return filepath.Join(config.Config.DataPath, file.ID.String())
}
