package repositories

import (
	"context"
	"fmt"
	"log/slog"

	"github.tomerab1/todo-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SongRepoIface interface {
	Create(ctx context.Context, song *models.Song) (*models.Song, error)
	FindAll(ctx context.Context) ([]*models.Song, error)
}

type SongRepoMongoImpl struct {
	coll   *mongo.Collection
	logger *slog.Logger
}

func NewSongRepoMongo(
	coll *mongo.Collection,
	logger *slog.Logger,
) SongRepoIface {
	return &SongRepoMongoImpl{
		coll:   coll,
		logger: logger,
	}
}

func (repo *SongRepoMongoImpl) Create(
	ctx context.Context,
	song *models.Song,
) (*models.Song, error) {
	if song.Id == "" {
		song.Id = primitive.NewObjectID().Hex()
	}

	_, err := repo.coll.InsertOne(ctx, song)
	if err != nil {
		return nil, fmt.Errorf("songRepo: %w: %v", ErrInsertFailed, err)
	}

	return song, nil
}

func (repo *SongRepoMongoImpl) FindAll(
	ctx context.Context,
) ([]*models.Song, error) {
	cursor, err := repo.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("songRepo: %w: %v", ErrFindAllFailed, err)
	}
	defer cursor.Close(ctx)

	var songs []*models.Song
	if err := cursor.All(ctx, &songs); err != nil {
		return nil, fmt.Errorf("songRepo: %w: %v", ErrFindAllFailed, err)
	}

	return songs, nil
}
