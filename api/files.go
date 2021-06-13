package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
)

func (a *api) allFiles(w http.ResponseWriter, r *http.Request) {
	var files []files.File

	if err := a.db.Find(&files).Error; err != nil {
		errorf(w, http.StatusInternalServerError, "find: %v", err)
		return
	}

	filesJSON := make([]map[string]interface{}, len(files))
	for i, file := range files {
		filesJSON[i] = map[string]interface{}{
			"id":   file.ID.String(),
			"kind": string(file.Kind),
			"time": file.Timestamp.String(),
			"live": string(file.LiveID),
			"dimensions": map[string]int{
				"width":  file.Width,
				"height": file.Height,
			},
			"endpoints": map[string]string{
				"data":  fmt.Sprintf("/files/data/%s%s", file.ID.String(), file.Extension),
				"thumb": fmt.Sprintf("/files/thumb/%s", file.ID.String()),
			},
		}
	}

	if err := json.NewEncoder(w).Encode(filesJSON); err != nil {
		errorf(w, http.StatusInternalServerError, "encode: %v", err)
		return
	}
}
