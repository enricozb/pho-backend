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
	_log.Debug("handling import")

	var importBody ImportBody

	if err := json.NewDecoder(req.Body).Decode(&importBody); err != nil {
		errorf(res, http.StatusBadRequest, "decode json: %v", err)
		return
	}

	_log.Infof("got body %v", importBody)
}
