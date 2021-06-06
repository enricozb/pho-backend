package workers

import (
	"fmt"
	"os"
	"time"

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

	var pathsToImport []paths.Path
	if err := w.db.Where("import_id = ?", importEntry.ID).Find(&pathsToImport).Error; err != nil {
		return fmt.Errorf("get paths: %v", err)
	}

	filesToImport := make([]files.File, len(pathsToImport))
	for i, path := range pathsToImport {
		metadata, err := validateMetadata(path.EXIFMetadata)
		if err != nil {
			return fmt.Errorf("extract metadata: %v", err)
		}

		filesToImport[i].ID = path.ID
		filesToImport[i].ImportID = path.ImportID
		filesToImport[i].Kind = path.Kind
		filesToImport[i].Timestamp = metadata.timestamp
		filesToImport[i].LiveID = metadata.liveID
		filesToImport[i].InitHash = path.InitHash
		filesToImport[i].Width = metadata.width
		filesToImport[i].Height = metadata.height
	}

	if err := w.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&filesToImport).Error; err != nil {
		return fmt.Errorf("insert files: %v", err)
	}

	if _, err := jobs.PushJob(w.db, importEntry.ID, jobs.JobConvert); err != nil {
		return fmt.Errorf("push job: %v", err)
	}

	return nil
}

type validatedEXIFMetadata struct {
	timestamp time.Time
	liveID    []byte
	width     int
	height    int
}

// validateMetadata checks for any missing EXIF data, setting defaults if necessary.
// TODO(enricozb): consider moving this to the EXIF worker.
func validateMetadata(exif paths.EXIFMetadata) (validatedEXIFMetadata, error) {
	var err error

	timestamp := time.Unix(exif.CreateDate, 0)
	if exif.CreateDate == 0 {
		timestamp, err = fallbackCreateDate(exif.Path)
		if err != nil {
			return validatedEXIFMetadata{}, fmt.Errorf("fallback create date: %v", err)
		}
	}

	if exif.Width == 0 || exif.Height == 0 {
		return validatedEXIFMetadata{}, fmt.Errorf("unable to get width/height: %v", err)
	}

	liveID := exif.MediaGroupUUID
	if liveID == "" {
		liveID = exif.ContentIdentifier
	}

	return validatedEXIFMetadata{
		timestamp: timestamp,
		liveID:    []byte(liveID),
		width:     exif.Width,
		height:    exif.Height,
	}, nil
}

func fallbackCreateDate(path string) (time.Time, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return time.Time{}, fmt.Errorf("stat: %v", err)
	}

	return fi.ModTime(), nil
}
