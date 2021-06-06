package api

import (
	"encoding/json"
	"net/http"
)

func (a *api) heartbeat(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		errorf(w, http.StatusInternalServerError, "encode: %v", err)
	}
}
