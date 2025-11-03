package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.tomerab1/todo-api/internal/app"
	"github.tomerab1/todo-api/internal/contracts"
)

func createSong(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto contracts.CreateSongDto
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			app.WriteErrorJSON(w, http.StatusBadRequest, fmt.Sprintf("invalid body: %v", err))
			return
		}

		resp, err := app.SongSvc.CreateSong(r.Context(), dto)
		if err != nil {
			app.WriteErrorJSON(w, http.StatusBadRequest, fmt.Sprintf("failed to create song: %v", err))
			return
		}

		app.WriteJSON(w, http.StatusCreated, resp)
	}
}

func getSongs(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		songs, err := app.SongSvc.GetAllSongs(r.Context())
		if err != nil {
			app.WriteErrorJSON(w, http.StatusBadRequest, fmt.Sprintf("failed to get songs: %v", err))
			return
		}

        app.WriteJSON(w, http.StatusOK, songs)
	}
}
