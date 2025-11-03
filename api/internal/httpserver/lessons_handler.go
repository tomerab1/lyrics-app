package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.tomerab1/todo-api/internal/app"
	"github.tomerab1/todo-api/internal/contracts"
	"github.tomerab1/todo-api/internal/services"
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
		var dto contracts.SubmitAnswerDto
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			app.WriteErrorJSON(w, http.StatusBadRequest, fmt.Sprintf("invalid body: %v", err))
			return
		}

		correct, err := app.LessonSvc.SubmitAnswer(r.Context(), dto.LessonId, dto.ItemIndex, dto.Type, dto.UserInput)
		if err != nil {
			if err == services.ErrDuplicateAnswer {
				app.WriteErrorJSON(w, http.StatusConflict, "duplicate submission")
				return
			}
			app.WriteErrorJSON(w, http.StatusBadRequest, fmt.Sprintf("failed to submit answer: %v", err))
			return
		}

		app.WriteJSON(w, http.StatusOK, contracts.SubmitAnswerResponse{Ok: true, Correct: correct})
	}
}

func lessonSummary(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// chi URL param
		lessonId := chi.URLParam(r, "lessonId")
		total, correct, wrong, acc, scheduled, err := app.LessonSvc.GetSummary(r.Context(), lessonId)
		if err != nil {
			app.WriteErrorJSON(w, http.StatusBadRequest, fmt.Sprintf("failed to get summary: %v", err))
			return
		}
		resp := contracts.LessonSummaryResponse{
			Total:                  total,
			Correct:                correct,
			Wrong:                  wrong,
			Accuracy:               acc,
			ScheduledForRepractice: scheduled,
		}

		app.WriteJSON(w, http.StatusOK, resp)
	}
}
