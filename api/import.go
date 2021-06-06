package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/workers/pkg/effects/workers"
)

type ImportBody struct {
	Opts jobs.ImportOptions `json:"opts"`
}

func (a *api) newImport(w http.ResponseWriter, r *http.Request) {
	_log.Debug("handling new import")

	var importBody ImportBody

	if err := json.NewDecoder(r.Body).Decode(&importBody); err != nil {
		errorf(w, http.StatusBadRequest, "decode json: %v", err)
		return
	}

	if importEntry, err := jobs.StartImport(a.db, importBody.Opts); err != nil {
		errorf(w, http.StatusInternalServerError, "start import: %v", err)
		return
	} else {
		if err := json.NewEncoder(w).Encode(map[string]string{
			"id": importEntry.String(),
		}); err != nil {
			errorf(w, http.StatusInternalServerError, "encode: %v", err)
			return
		}
	}
}

func (a *api) cleanupImport(w http.ResponseWriter, r *http.Request) {
	_log.Debug("handling cleanup import")

	importID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		errorf(w, http.StatusBadRequest, "malformed import id: %v", err)
		return
	}

	if _, err := jobs.PushJobWithArgs(a.db, importID, jobs.JobCleanup, workers.CleanupWorkerArgs{Full: true}); err != nil {
		errorf(w, http.StatusInternalServerError, "start import: %v", err)
	}
}

func (a *api) importStatus(w http.ResponseWriter, r *http.Request) {
	_log.Debug("handling import status")

	importID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		errorf(w, http.StatusBadRequest, "malformed import id: %v", err)
		return
	}

	importEntry := jobs.Import{ID: importID}

	if err := a.db.First(&importEntry).Error; err != nil {
		errorf(w, http.StatusInternalServerError, "get import: %v", err)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{
		"id":         importEntry.ID.String(),
		"status":     string(importEntry.Status),
		"updated_at": importEntry.UpdatedAt.String(),
	}); err != nil {
		errorf(w, http.StatusInternalServerError, "encode: %v", err)
		return
	}
}
