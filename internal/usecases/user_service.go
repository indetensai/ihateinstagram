package usecases

import (
	"context"
	"phota/internal/entities"
	"phota/internal/usecases/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type userService struct {
	db *pgx.Conn
}

func (u *userService) CheckSession(session_id string, ctx context.Context) (*uuid.UUID, error) {
	if session_id == "" {
		return nil, entities.ErrEmptySession
	}
	user_repository := repository.NewUserRepository(u.db)
	user_id, err := user_repository.CheckSession(session_id, ctx)
	if err != nil {
		return nil, err
	}
	return user_id, nil
}

func (u *userService) DeleteSession(session_id string, ctx context.Context) error {
	user_repository := repository.NewUserRepository(u.db)
	err := user_repository.DeleteSession(session_id, ctx)
	return err
}

func (u *userService) Login(username string, password string, ctx context.Context) (*string, error) {
	user_repository := repository.NewUserRepository(u.db)
	session_id, err := user_repository.Login(username, password, ctx)
	if err != nil {
		return nil, err
	}

	return session_id, nil
}

func (u *userService) Register(username string, password string, ctx context.Context) error {
	user_repository := repository.NewUserRepository(u.db)
	err := user_repository.Register(username, password, ctx)
	return err
}
