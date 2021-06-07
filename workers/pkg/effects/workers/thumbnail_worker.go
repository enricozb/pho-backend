package workers

import (
	"fmt"
	"path/filepath"

	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/config"
	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/thumbnail"
	"github.com/enricozb/pho/workers/pkg/lib/worker"
)

type ThumbnailWorker struct {
	db *gorm.DB
}

var _ worker.Worker = &ThumbnailWorker{}

func NewThumbnailWorker(db *gorm.DB) *ThumbnailWorker {
	return &ThumbnailWorker{db: db}
}

func (w *ThumbnailWorker) Work(job jobs.Job) error {
	importEntry := jobs.Import{}
	if err := w.db.First(&importEntry, job.ImportID).Error; err != nil {
		return fmt.Errorf("find import: %v", err)
	}

	if err := importEntry.SetStatus(w.db, jobs.ImportStatusThumbnail); err != nil {
		return fmt.Errorf("set import status: %v", err)
	}

	var filesToImport []files.File
	if err := w.db.Where("import_id = ?", importEntry.ID, files.ImageKind).Find(&filesToImport).Error; err != nil {
		return fmt.Errorf("get image files: %v", err)
	}

	// generate the thumbnails
	thumbs := thumbnail.NewThumbnailGenerator()

	for _, file := range filesToImport {
		src := thumbnailSrcPath(file)
		dst := thumbnailDestPath(file)
		if err := thumbs.Thumbnail(src, dst, file.Kind); err != nil {
			return fmt.Errorf("thumbnail: %v", err)
		}
	}

	if err := thumbs.Finish(); err != nil {
		return fmt.Errorf("finish: %v", err)
	}

	if _, err := jobs.PushJobWithArgs(w.db, importEntry.ID, jobs.JobCleanup, CleanupWorkerArgs{Full: false}); err != nil {
		return fmt.Errorf("push job: %v", err)
	}

	return nil
}

// thumbnailSrcPath returns the source path of a file after its conversion, which will be used as the source for the thumbnail.
func thumbnailSrcPath(file files.File) string {
	return filepath.Join(config.Config.MediaDir, file.ID.String()+file.Extension)
}

// thumbnailDestPath returns the destination path after the creation of a thumbnail, with the thumbnail extension.
func thumbnailDestPath(file files.File) string {
	return filepath.Join(config.Config.ThumbDir, file.ID.String())
}
