package ports

import (
	"codebase-app/internal/module/user/entity"
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req *entity.CreateUserRequest) (*entity.CreateUserResponse, error)
}

type UserService interface {
	CreateUser(ctx context.Context, req *entity.CreateUserRequest) (*entity.CreateUserResponse, error)
}
