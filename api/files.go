package api

import (
	"encoding/json"
	"net/http"

	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
)

func (a *api) allFiles(w http.ResponseWriter, r *http.Request) {
	var files []files.File

	if err := a.db.Find(&files).Error; err != nil {
		errorf(w, http.StatusInternalServerError, "find: %v", err)
		return
	}

	filesJSON := make(map[string]interface{})
	for _, file := range files {
		filesJSON[file.ID.String()] = map[string]string{
			"kind": string(file.Kind),
			"time": file.Timestamp.UTC().String(),
			"live": string(file.LiveID),
		}
	}

	if err := json.NewEncoder(w).Encode(filesJSON); err != nil {
		errorf(w, http.StatusInternalServerError, "encode: %v", err)
		return
	}
}
