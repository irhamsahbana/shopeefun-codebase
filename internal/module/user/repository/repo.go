package repository

import (
	"codebase-app/internal/module/user/entity"
	"codebase-app/internal/module/user/ports"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.UserRepository = &userRepository{}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, req *entity.CreateUserRequest) (*entity.CreateUserResponse, error) {
	var res = new(entity.CreateUserResponse)
	// your implementation here
	query := `
		INSERT INTO
			users (role_id, name, email, password)
		VALUES (?, ?, ?, 'undefined')
			RETURNING id`

	err := r.db.QueryRowContext(ctx, r.db.Rebind(query), req.RoleId, req.Name, req.Email).Scan(&res.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::CreateUser - Failed to create user")
		return nil, err
	}

	return res, nil
}
