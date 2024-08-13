package service

import (
	"codebase-app/internal/module/user/entity"
	"codebase-app/internal/module/user/ports"
	"context"
)

var _ ports.UserService = &userService{}

type userService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *userService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *entity.CreateUserRequest) (*entity.CreateUserResponse, error) {
	return s.repo.CreateUser(ctx, req)
}
