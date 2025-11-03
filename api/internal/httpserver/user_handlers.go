package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.tomerab1/todo-api/internal/app"
	"github.tomerab1/todo-api/internal/contracts"
)

func createUser(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto contracts.CreateUserDto
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			app.WriteErrorJSON(w, http.StatusBadRequest, fmt.Sprintf("invalid body: %v", err))
			return
		}

		resp, err := app.UserSvc.CreateUser(r.Context(), dto)
		if err != nil {
			app.WriteErrorJSON(w, http.StatusBadRequest, fmt.Sprintf("failed to create user: %v", err))
			return
		}

		app.WriteJSON(w, http.StatusCreated, resp)
	}
}

func getUsers(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := app.UserSvc.GetAllUsers(r.Context())
		if err != nil {
			app.WriteErrorJSON(w, http.StatusBadRequest, fmt.Sprintf("failed to get users: %v", err))
			return
		}

        app.WriteJSON(w, http.StatusOK, users)
	}
}
