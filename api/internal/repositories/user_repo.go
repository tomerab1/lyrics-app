package repositories

import (
	"context"
	"fmt"
	"log/slog"

	"github.tomerab1/todo-api/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepoIface interface {
	Create(ctx context.Context, user *models.User) (string, error)
	FindAll(ctx context.Context) ([]*models.User, error)
}

type UserRepoMongoImpl struct {
	coll   *mongo.Collection
	logger *slog.Logger
}

func NewUserRepoMongo(
	coll *mongo.Collection,
	logger *slog.Logger,
) UserRepoIface {
	return &UserRepoMongoImpl{
		coll:   coll,
		logger: logger,
	}
}

func (repo *UserRepoMongoImpl) Create(
	ctx context.Context,
	user *models.User,
) (string, error) {
	if user.Id == "" {
		user.Id = primitive.NewObjectID().Hex()
	}

	_, err := repo.coll.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("userRepo: %w: %v", ErrInsertFailed, err)
	}

	return user.Id, nil
}

func (repo *UserRepoMongoImpl) FindAll(
	ctx context.Context,
) ([]*models.User, error) {
	cursor, err := repo.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("userRepo: %w: %v", ErrFindAllFailed, err)
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("userRepo: %w: %v", ErrFindAllFailed, err)
	}

	return users, nil
}
