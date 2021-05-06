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

	// grab exif data
	var filepaths []string
	for _, path := range pathEntries {
		filepaths = append(filepaths, path.Path)
	}

	exifData, err := getEXIFData(importEntry.ID, filepaths)
	if err != nil {
		return fmt.Errorf("get exif data: %v", err)
	}

	for _, path := range pathEntries {
		exif, ok := exifData[path.Path]
		if !ok {
			return fmt.Errorf("missing exif data from path: %s", path.Path)
		}

		path.EXIFMetadata = exif
		if err := w.db.Save(&path).Error; err != nil {
			return fmt.Errorf("update: %v", err)
		}
	}

	return nil
}

// getEXIFData returns `EXIFMetadata` for each of path in `paths`, returning a mapping from `path.Path` to `EXIFMetadata`.
func getEXIFData(importID jobs.ImportID, filepaths []string) (map[string]paths.EXIFMetadata, error) {
	// write file paths to temporary file
	tmp, err := ioutil.TempFile("", "pho-import-files-"+importID.String())
	for _, path := range filepaths {
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

	var exifMetadatas []paths.EXIFMetadata
	if err := json.Unmarshal(data, &exifMetadatas); err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}

	exifMap := map[string]paths.EXIFMetadata{}
	for _, exif := range exifMetadatas {
		exifMap[exif.Path] = exif
	}

	return exifMap, nil
}
