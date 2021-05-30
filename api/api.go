package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/lib/logs"
)

var _log = logs.MustGetLogger("scheduler")

type api struct {
	db *gorm.DB
}

func NewAPI(db *gorm.DB) *api {
	return &api{db: db}
}

func (a *api) Run() error {
	_log.Debug("running api...")

	r := mux.NewRouter()

	r.HandleFunc("/import", handleImport).Methods("POST")

	http.Handle("/", r)

	return http.ListenAndServe(":4000", r)
}
