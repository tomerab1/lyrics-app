package services

import (
	"context"
	"log/slog"

	"github.tomerab1/todo-api/internal/contracts"
	"github.tomerab1/todo-api/internal/models"
	"github.tomerab1/todo-api/internal/repositories"
	"github.tomerab1/todo-api/internal/utils"
)

type SongService struct {
	songRepo repositories.SongRepoIface
	logger   *slog.Logger
}

func NewSongService(
	repo repositories.SongRepoIface,
	logger *slog.Logger,
) *SongService {
	return &SongService{
		songRepo: repo,
		logger:   logger,
	}
}

func (svc *SongService) CreateSong(
	ctx context.Context,
	createSongDto contracts.CreateSongDto,
) (*contracts.CreateSongsReponse, error) {
	song, err := svc.songRepo.Create(
		ctx,
		&models.Song{
			Title:  createSongDto.Title,
			Artist: createSongDto.Artist,
			Lyrics: utils.LyricsToSlices(createSongDto.Lyrics),
		},
	)

	if err != nil {
		return nil, err
	}

	return &contracts.CreateSongsReponse{
		Id:        song.Id,
		LineCount: len(song.Lyrics),
	}, nil
}

func (svc *SongService) GetAllSongs(
	ctx context.Context,
) ([]contracts.GetSongResponse, error) {
	songs, err := svc.songRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]contracts.GetSongResponse, 0)
	for _, song := range songs {
		resp = append(resp, contracts.GetSongResponse{
			Id:    song.Id,
			Title: song.Title,
		})
	}

	return resp, nil
}
