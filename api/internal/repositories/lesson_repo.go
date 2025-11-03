// repositories/lesson_repo.go
package repositories

import (
    "context"
    "fmt"
    "log/slog"
    "time"

    "github.tomerab1/todo-api/internal/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/v2/mongo"
)

type LessonRepoIface interface {
    Create(ctx context.Context, userId string, lesson *models.Lesson) (*models.Lesson, error)
    GetById(ctx context.Context, id string) (*models.Lesson, error)
    AddAnswer(ctx context.Context, lessonId string, ans models.LessonAnswer) error
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

func (repo *LessonRepoMongoDb) GetById(
    ctx context.Context,
    id string,
) (*models.Lesson, error) {
    var out models.Lesson
    err := repo.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&out)
    if err != nil {
        return nil, fmt.Errorf("lessonRepo: find failed: %w", err)
    }
    return &out, nil
}

func (repo *LessonRepoMongoDb) AddAnswer(
    ctx context.Context,
    lessonId string,
    ans models.LessonAnswer,
) error {
    // Ensure answers field is an array (convert null -> [])
    _, _ = repo.coll.UpdateOne(ctx,
        bson.M{"_id": lessonId, "answers": bson.M{"$type": "null"}},
        bson.M{"$set": bson.M{"answers": bson.A{}}},
    )
    filter := bson.M{"_id": lessonId, "answers.item_index": bson.M{"$ne": ans.ItemIndex}}
    update := bson.M{"$push": bson.M{"answers": ans}}
    res, err := repo.coll.UpdateOne(ctx, filter, update)
    if err != nil {
        return fmt.Errorf("lessonRepo: add answer failed: %w", err)
    }
    if res.MatchedCount == 0 {
        return fmt.Errorf("lessonRepo: duplicate answer or lesson not found")
    }
    return nil
}
