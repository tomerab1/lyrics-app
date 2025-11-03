package app

import (
	"log/slog"

	"github.tomerab1/todo-api/internal/repositories"
	"github.tomerab1/todo-api/internal/services"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Application struct {
	db        *mongo.Client
	logger    *slog.Logger
	UserSvc   *services.UserService
	SongSvc   *services.SongService
	LessonSvc *services.LessonService
}

func New(logger *slog.Logger, dbConnString string) (*Application, error) {
	dbConn, err := mongo.Connect(options.Client().ApplyURI(dbConnString))
	if err != nil {
		return nil, err
	}

	userRepoLogger := slog.New(logger.Handler()).With("repo", "user")
	userSvcLogger := slog.New(logger.Handler()).With("service", "user")

	songRepoLogger := slog.New(logger.Handler()).With("repo", "songs")
	songSvcLogger := slog.New(logger.Handler()).With("service", "songs")

	lessonsRepoLogger := slog.New(logger.Handler()).With("repo", "lessons")
	lessonsSvcLogger := slog.New(logger.Handler()).With("service", "lessons")

	userRepo := repositories.NewUserRepoMongo(dbConn.Database("lyrics-app").Collection("users"), userRepoLogger)
	userSvc := services.NewUserService(userRepo, userSvcLogger)

	songsRepo := repositories.NewSongRepoMongo(dbConn.Database("lyrics-app").Collection("songs"), songRepoLogger)
	songsSvc := services.NewSongService(songsRepo, songSvcLogger)

	lessonRepo := repositories.NewLessonRepo(dbConn.Database("lyrics-app").Collection("lessons"), lessonsRepoLogger)
	lessonSvc := services.NewLessonService(songsRepo, lessonRepo, lessonsSvcLogger)

	return &Application{
		db:        dbConn,
		UserSvc:   userSvc,
		SongSvc:   songsSvc,
		LessonSvc: lessonSvc,
	}, nil
}
