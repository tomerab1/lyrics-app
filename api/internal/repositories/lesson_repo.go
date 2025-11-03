// repositories/lesson_repo.go
package repositories

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.tomerab1/todo-api/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type LessonRepoIface interface {
	Create(ctx context.Context, userId string, lesson *models.Lesson) (*models.Lesson, error)
}

type LessonRepoMongoDb struct {
	coll   *mongo.Collection
	logger *slog.Logger
}

func NewLessonRepo(coll *mongo.Collection, logger *slog.Logger) LessonRepoIface {
	return &LessonRepoMongoDb{coll: coll, logger: logger}
}

func (repo *LessonRepoMongoDb) Create(
	ctx context.Context,
	userId string,
	lesson *models.Lesson,
) (*models.Lesson, error) {
	if lesson == nil {
		return nil, fmt.Errorf("lessonRepo: nil lesson")
	}
	if lesson.Id == "" {
		lesson.Id = primitive.NewObjectID().Hex()
	}
	lesson.UserId = userId
	if lesson.CreatedAt.IsZero() {
		lesson.CreatedAt = time.Now().UTC()
	}
	_, err := repo.coll.InsertOne(ctx, lesson)
	if err != nil {
		return nil, fmt.Errorf("lessonRepo: insert failed: %w", err)
	}
	return lesson, nil
}
