package workers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"gorm.io/gorm"

	"github.com/op/go-logging"

	"github.com/enricozb/pho/shared/pkg/effects/daos/exif"
	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/daos/paths"
	"github.com/enricozb/pho/shared/pkg/lib/logs"
	"github.com/enricozb/pho/workers/pkg/lib/worker"
)

type exifWorker struct {
	db  *gorm.DB
	log *logging.Logger
}

var _ worker.Worker = &exifWorker{}

func NewEXIFWorker(db *gorm.DB) *exifWorker {
	return &exifWorker{db: db, log: logs.MustGetLogger("exif worker")}
}

func (w *exifWorker) Work(job jobs.Job) error {
	importEntry := jobs.Import{}
	if err := w.db.First(&importEntry, job.ImportID).Error; err != nil {
		return fmt.Errorf("find import: %v", err)
	}

	pathEntries, err := paths.PathsInPipeline(w.db, importEntry.ID)
	if err != nil {
		return fmt.Errorf("get paths: %v", err)
	}

	// grab exif data
	var filepaths []string
	for _, path := range pathEntries {
		filepaths = append(filepaths, path.Path)
	}

	exifData, err := w.getEXIFData(importEntry.ID, filepaths)
	if err != nil {
		return fmt.Errorf("get exif data: %v", err)
	}

	for _, path := range pathEntries {
		var ok bool

		path.EXIFMetadata, ok = exifData[path.Path]
		if !ok {
			return fmt.Errorf("missing exif data from path: %s", path.Path)
		}

		if err := w.db.Model(&path).Select("exif_metadata").Updates(&path).Error; err != nil {
			return fmt.Errorf("update: %v", err)
		}
	}

	return nil
}

type exifExtraMetadata struct {
	exif.EXIFMetadata

	MediaGroupUUID    string `json:"MediaGroupUUID"`
	ImageUniqueID     string `json:"ImageUniqueID"`
	ContentIdentifier string `json:"ContentIdentifier"`
	Orientation       int    `json:"Orientation"`

	CreateDate       string `json:"CreateDate"`
	CreationDate     string `json:"CreationDate"`
	DateTimeOriginal string `json:"DateTimeOriginal"`
}

// normalize modifies any fields before saving to the DB. Some operations include picking the best timestamp field and fixing the orientation for HEIC images.
func (e *exifExtraMetadata) normalize() (_ exif.EXIFMetadata, err error) {
	// need to flip the width and height if depending on orientation metadata
	// see: https://sirv.com/help/articles/rotate-photos-to-be-upright/
	if e.Orientation >= 5 {
		e.Width, e.Height = e.Height, e.Width
	}

	// picking live id
	if e.MediaGroupUUID != "" {
		e.LiveID = []byte(e.MediaGroupUUID)
	} else if e.ContentIdentifier != "" {
		e.LiveID = []byte(e.ContentIdentifier)
	}

	// picking the best timestamp representation...
	if e.CreationDate != "" {
		e.Timestamp = e.CreationDate
	} else if e.DateTimeOriginal != "" {
		e.Timestamp = e.DateTimeOriginal
	} else if e.CreateDate != "" {
		e.Timestamp = e.CreateDate
	} else {
		e.Timestamp, err = e.filesystemTimestamp()
	}

	if e.Width == 0 || e.Height == 0 {
		return exif.EXIFMetadata{}, fmt.Errorf("missing width/height: %s", e.Path)
	}

	return e.EXIFMetadata, err
}

func (e *exifExtraMetadata) filesystemTimestamp() (string, error) {
	fi, err := os.Stat(e.Path)
	if err != nil {
		return "", fmt.Errorf("stat: %v", err)
	}

	return fi.ModTime().Format("2006:01:02 15:04:05-07:00"), nil
}

// getEXIFData returns `EXIFMetadata` for each of path in `paths`, returning a mapping from `path.Path` to `EXIFMetadata`.
func (w *exifWorker) getEXIFData(importID jobs.ImportID, filepaths []string) (map[string]exif.EXIFMetadata, error) {
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

		// timestamp information
		"-CreationDate",
		"-DateTimeOriginal",
		"-CreateDate",

		// additional exif metadata
		"-MediaGroupUUID",
		"-ImageUniqueID",
		"-ContentIdentifier",
		"-ImageWidth",
		"-ImageHeight",
		"-Orientation#",
	)

	data, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("run exiftool [%v]: stderr: %s", err, data)
	}

	var exifMetadatas []exifExtraMetadata
	if err := json.Unmarshal(data, &exifMetadatas); err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}

	exifMap := map[string]exif.EXIFMetadata{}
	for _, exif := range exifMetadatas {
		exifMap[exif.Path], err = exif.normalize()
		if err != nil {
			return nil, fmt.Errorf("normalize: %v", err)
		}
	}

	return exifMap, nil
}
