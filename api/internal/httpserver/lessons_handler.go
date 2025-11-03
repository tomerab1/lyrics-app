package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.tomerab1/todo-api/internal/app"
	"github.tomerab1/todo-api/internal/contracts"
)

func createLesson(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req contracts.CreateLessonDto
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			app.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse lesson: %v", err))
			return
		}
		out, err := app.LessonSvc.CreateLesson(r.Context(), contracts.CreateLessonDto{UserId: req.UserId})
		if err != nil {
			app.WriteErrorJSON(w, http.StatusBadRequest, fmt.Sprintf("failed to create lesson: %v", err))
			return
		}

		app.WriteJSON(w, http.StatusCreated, out)
	}
}

func submitAnswer(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func lessonSummary(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
