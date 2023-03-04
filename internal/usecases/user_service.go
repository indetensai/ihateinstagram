package usecases

import (
	"context"
	"phota/internal/entities"

	"github.com/google/uuid"
)

type userService struct {
	repo entities.UserRepository
}

func NewUserService(u entities.UserRepository) entities.UserService {
	return &userService{u}
}

func (u *userService) CheckSession(session_id string, ctx context.Context) (*uuid.UUID, error) {
	if session_id == "" {
		return nil, entities.ErrEmptySession
	}
	user_id, err := u.repo.CheckSession(session_id, ctx)
	if err != nil {
		return nil, err
	}
	return user_id, nil
}

func (u *userService) DeleteSession(session_id string, ctx context.Context) error {
	err := u.repo.DeleteSession(session_id, ctx)
	return err
}

func (u *userService) Login(username string, password string, ctx context.Context) (*string, error) {
	session_id, err := u.repo.Login(username, password, ctx)
	if err != nil {
		return nil, err
	}

	return session_id, nil
}

func (u *userService) Register(username string, password string, ctx context.Context) error {
	err := u.repo.Register(username, password, ctx)
	return err
}
