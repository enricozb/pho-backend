package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/enricozb/pho/shared/pkg/effects/daos/albums"
	"github.com/enricozb/pho/shared/pkg/effects/daos/files"
)

func (a *api) allAlbums(w http.ResponseWriter, r *http.Request) {
	var albums []albums.Album

	if err := a.db.Find(&albums).Error; err != nil {
		errorf(w, http.StatusInternalServerError, "find: %v", err)
		return
	}

	var albumsJSON []map[string]string

	for _, album := range albums {
		albumsJSON = append(albumsJSON, map[string]string{
			"id":   album.ID.String(),
			"name": album.Name,
		})
	}

	if err := json.NewEncoder(w).Encode(albumsJSON); err != nil {
		errorf(w, http.StatusInternalServerError, "encode: %v", err)
		return
	}
}

type UpdateAlbumBody struct {
	Name  string   `json:"name"`
	Files []string `json:"files"`
}

func (a *api) newAlbum(w http.ResponseWriter, r *http.Request) {
	var albumBody UpdateAlbumBody

	if err := json.NewDecoder(r.Body).Decode(&albumBody); err != nil {
		errorf(w, http.StatusBadRequest, "decode json: %v", err)
		return
	}

	files := []files.File{}
	if err := a.db.Find(&files, albumBody.Files).Error; err != nil {
		errorf(w, http.StatusBadRequest, "bad file id: %v", err)
		return
	}

	album := albums.Album{Name: albumBody.Name, Files: files}
	if err := a.db.Create(&album).Error; err != nil {
		errorf(w, http.StatusBadRequest, "new album: %v", err)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"id": album.ID.String()}); err != nil {
		errorf(w, http.StatusInternalServerError, "encode: %v", err)
		return
	}
}

func (a *api) albumCover(w http.ResponseWriter, r *http.Request) {
	albumID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		errorf(w, http.StatusBadRequest, "malformed album id: %v", err)
		return
	}

	album := albums.Album{ID: albumID}
	if err := a.db.Preload("Files").First(&album).Error; err != nil {
		errorf(w, http.StatusBadRequest, "get album: %v", err)
		return
	}

	coverJSON := map[string]string{}
	if len(album.Files) > 0 {
		coverJSON["cover"] = fmt.Sprintf("/files/thumb/%s", album.Files[0].ID.String())
	}

	if err := json.NewEncoder(w).Encode(coverJSON); err != nil {
		errorf(w, http.StatusBadRequest, "encode: %v", err)
		return
	}
}

func (a *api) albumData(w http.ResponseWriter, r *http.Request) {
	albumID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		errorf(w, http.StatusBadRequest, "malformed album id: %v", err)
		return
	}

	album := albums.Album{ID: albumID}
	if err := a.db.Preload("Files").First(&album).Error; err != nil {
		errorf(w, http.StatusBadRequest, "get album: %v", err)
		return
	}

	albumJSON := map[string]interface{}{
		"id":    albumID.String(),
		"name":  album.Name,
		"files": a.filesJSON(album.Files),
		// TODO(enricozb) add parent and children
	}

	if err := json.NewEncoder(w).Encode(albumJSON); err != nil {
		errorf(w, http.StatusBadRequest, "encode: %v", err)
		return
	}
}

func (a *api) updateAlbum(w http.ResponseWriter, r *http.Request) {
	var albumBody UpdateAlbumBody
	if err := json.NewDecoder(r.Body).Decode(&albumBody); err != nil {
		errorf(w, http.StatusBadRequest, "decode json: %v", err)
		return
	}

	albumID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		errorf(w, http.StatusBadRequest, "malformed album id: %v", err)
		return
	}
	album := albums.Album{ID: albumID}

	// rename
	if len(albumBody.Name) > 0 {
		if err := a.db.Debug().Model(&album).Update("name", albumBody.Name).Error; err != nil {
			errorf(w, http.StatusBadRequest, "save album: %v", err)
			return
		}
	}

	// add files
	// TODO(enricozb) handle removing files
	if len(albumBody.Files) == 0 {
		files := []files.File{}
		if err := a.db.Debug().Find(&files, albumBody.Files).Error; err != nil {
			errorf(w, http.StatusBadRequest, "bad file id: %v", err)
			return
		}

		if err := a.db.Model(&album).Association("Files").Append(files); err != nil {
			errorf(w, http.StatusBadRequest, "append files: %v", err)
			return
		}
	}
}
