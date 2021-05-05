package workers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
	"github.com/enricozb/pho/workers/pkg/lib/worker"
)

type EXIFData struct {
	Path              string `json:"SourceFile"`
	CreateDate        int64  `json:"CreateDate"`
	MediaGroupUUID    string `json:"MediaGroupUUID"`
	ImageUniqueID     string `json:"ImageUniqueID"`
	ContentIdentifier string `json:"ContentIdentifier"`
}

type exifWorker struct {
	db *gorm.DB
}

var _ worker.Worker = &exifWorker{}

func NewEXIFWorker(db *gorm.DB) *exifWorker {
	return &exifWorker{db: db}
}

func (w *exifWorker) Work(job jobs.Job) error {
	importEntry := jobs.Import{}
	if err := w.db.Find(&importEntry, job.ImportID).Error; err != nil {
		return fmt.Errorf("find import: %v", err)
	}

	var pathEntries []paths.Path
	if err := w.db.Where("import_id = ?", importEntry.ID).Find(&pathEntries).Error; err != nil {
		return fmt.Errorf("get paths: %v", err)
	}

	// maps file path to path uuids
	pathIDs := map[string]string{}
	var osPaths []string
	for _, path := range pathEntries {
		osPaths = append(osPaths, path.Path)
		pathIDs[path.Path] = path.ID.String()
	}

	// grab exif data
	exifData, err := getEXIFData(importEntry.ID, osPaths)
	if err != nil {
		return fmt.Errorf("get exif data: %v", err)
	}

	// insert exif data into `paths` table
	for _, exif := range exifData {
		pathID, ok := pathIDs[exif.Path]
		if !ok {
			return fmt.Errorf("got exif data for non-existing path: %s", exif.Path)
		}

		exifJSON, err := json.Marshal(exif)
		if err != nil {
			return fmt.Errorf("marshal: %v", err)
		}

		if err := w.db.Model(&paths.Path{}).Where("id = ?", pathID).Update("exif_metadata", exifJSON).Error; err != nil {
			return fmt.Errorf("update: %v", err)
		}
	}

	return nil
}

func getEXIFData(importID jobs.ImportID, paths []string) (exifData []EXIFData, err error) {
	// write file paths to temporary file
	tmp, err := ioutil.TempFile("", "pho-import-files-"+importID.String())
	for _, path := range paths {
		tmp.Write(append([]byte(path), '\n'))
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	cmd := exec.Command(
		"exiftool",
		"-@", tmp.Name(),
		"-json",
		// desired exif metadata
		"-CreateDate", "-DateFormat", "%s",
		"-MediaGroupUUID",
		"-ImageUniqueID",
		"-ContentIdentifier",
	)

	data, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("run exiftool: %v", err)
	}

	return exifData, json.Unmarshal(data, &exifData)
}
