package app

import (
	"encoding/json"
	"net/http"
)

func (a *Application) WriteJSON(
	w http.ResponseWriter,
	status int,
	data any,
) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"data": data,
	})
}

func (a *Application) WriteErrorJSON(
	w http.ResponseWriter,
	status int,
	errMsg string,
) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"error": errMsg,
	})
}
