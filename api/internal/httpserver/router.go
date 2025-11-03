package httpserver

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.tomerab1/todo-api/internal/app"
)

func New(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(5 * time.Second))
	r.Use(commonHeadersMiddleware)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	api := chi.NewRouter()

	api.Route("/users", func(r chi.Router) {
		r.Post("/", createUser(app))
		r.Get("/", getUsers(app))
	})

	api.Route("/songs", func(r chi.Router) {
		r.Post("/", createSong(app))
		r.Get("/", getSongs(app))
	})

	api.Post("/lessons", createLesson(app))
	api.Post("/answers", submitAnswer(app))
	api.Get("/lessons/{lessonId}/summary", lessonSummary(app))

	r.Mount("/api", api)
	return r
}
