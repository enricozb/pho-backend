package api

import (
	"encoding/json"
	"net/http"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
)

type ImportBody struct {
	Opts jobs.ImportOptions `json:"opts"`
}

func handleImport(res http.ResponseWriter, req *http.Request) {
	_log.Info("handling import")

	var importBody ImportBody

	if err := json.NewDecoder(req.Body).Decode(&importBody); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}
