package entities

import (
	"context"

	"github.com/google/uuid"
)

type AppUser struct {
	UserID   uuid.UUID
	Username string
}

type UserCredentials struct {
	UserID   uuid.UUID
	Username string
	Password string
}
type UserService interface {
	CheckSession(session_id string, ctx context.Context) (*uuid.UUID, error)
	DeleteSession(session_id string, ctx context.Context) error
	Login(username string, password string, ctx context.Context) (*string, error)
	Register(username string, password string, ctx context.Context) error
}

type UserRepository interface {
	GetUserBySession(session_id string, ctx context.Context) (*uuid.UUID, error)
	DeleteSession(session_id string, ctx context.Context) error
	ValidateUserCredentials(username string, password string, session_id string, ctx context.Context) error
	CreateUser(username string, password string, ctx context.Context) error
	GetUserByID(user_id uuid.UUID, ctx context.Context) (*AppUser, error)
}
