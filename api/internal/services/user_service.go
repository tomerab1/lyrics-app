package services

import (
	"context"
	"log/slog"

	"github.tomerab1/todo-api/internal/contracts"
	"github.tomerab1/todo-api/internal/models"
	"github.tomerab1/todo-api/internal/repositories"
)

type UserService struct {
	userRepo repositories.UserRepoIface
	logger   *slog.Logger
}

func NewUserService(
	repo repositories.UserRepoIface,
	logger *slog.Logger,
) *UserService {
	return &UserService{
		userRepo: repo,
		logger:   logger,
	}
}

func (svc *UserService) CreateUser(
	ctx context.Context,
	createUserDto contracts.CreateUserDto,
) (string, error) {
	return svc.userRepo.Create(ctx, &models.User{
		Name: createUserDto.Name,
	})
}

func (svc *UserService) GetAllUsers(
	ctx context.Context,
) ([]contracts.GetUserResponse, error) {
	users, err := svc.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]contracts.GetUserResponse, 0)
	for _, user := range users {
		resp = append(resp, contracts.GetUserResponse{
			Id:   user.Id,
			Name: user.Name,
		})
	}

	return resp, nil
}
