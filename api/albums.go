package api

import (
	"encoding/json"
	"net/http"

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

type NewAlbumBody struct {
	Name  string   `json:"name"`
	Files []string `json:"files"`
}

func (a *api) newAlbum(w http.ResponseWriter, r *http.Request) {
	var albumBody NewAlbumBody

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
