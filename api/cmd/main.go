package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.tomerab1/todo-api/internal/app"
	"github.tomerab1/todo-api/internal/httpserver"
)

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	_ = godotenv.Load()

	mongoURI := getenv("MONGODB_URI", "")
	if mongoURI == "" {
		logger.Warn("MONGODB_URI is empty")
	}

	port := getenv("PORT", "8080")
	addr := ":" + port

	app, err := app.New(logger, mongoURI)
	if err != nil {
		logger.Error("failed to create app", "err", err)
		os.Exit(1)
	}

	handler := httpserver.New(app)

	mux := http.NewServeMux()
	mux.Handle("/", handler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	logger.Info("server starting", "addr", addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen and serve failed", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", "err", err)
	}
	logger.Info("server stopped")
}
