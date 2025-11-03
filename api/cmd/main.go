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

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	if err := godotenv.Load(); err != nil {
		logger.Error("godotenv load failed", "err", err)
		panic(".env load failed")
	}

	var (
		SERVER_ADDR = os.Getenv("SERVER_ADDR")
		MONGO_ADDR  = os.Getenv("MONGO_ADDR")
	)

	app, err := app.New(logger, MONGO_ADDR)
	if err != nil {
		logger.Error("failed to create app", "err", err)
		panic("failed to create app")
	}

	srv := http.Server{
		Addr:    SERVER_ADDR,
		Handler: httpserver.New(app),
	}
	logger.Info("server is running", "addr", srv.Addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Could not listen and serve", "addr", srv.Addr, "err", err.Error())
			panic("failed to listen and serve")
		}
	}()

	<-quit
	logger.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Info("server shutdown failed", "err", err.Error())
	}

	logger.Info("server gracefully stopped.")
}
