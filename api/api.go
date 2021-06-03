package api

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/config"
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

	r.HandleFunc("/heartbeat", a.heartbeat).Methods("GET")

	r.HandleFunc("/import/new", a.newImport).Methods("POST")
	r.HandleFunc("/import/{id:[-0-9a-f]+}/status", a.importStatus).Methods("GET")

	r.HandleFunc("/files/all", a.allFiles).Methods("GET")
	r.PathPrefix("/files/data/").Handler(http.StripPrefix("/files/data/", http.FileServer(http.Dir(config.Config.MediaDir))))

	http.Handle("/", r)

	return http.ListenAndServe(":4000", handlers.CORS()(r))
}
