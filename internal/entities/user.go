package entities

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AppUser struct {
	UserID   uuid.UUID
	Name     string
	Password string
}

type Session struct {
	SessionID uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
}

type UserHandler interface {
	RegisterHandler(c *fiber.Ctx) error
	LoginHandler(c *fiber.Ctx) error
	DeleteSessionHandler(c *fiber.Ctx) error
}

type UserService interface {
	CheckSession(session_id string, ctx context.Context) (*uuid.UUID, error)
	DeleteSession(session_id string, ctx context.Context) error
	Login(username string, password string, ctx context.Context) (*string, error)
	Register(username string, password string, ctx context.Context) error
}

type UserRepository interface {
	CheckSession(session_id string, ctx context.Context) (*uuid.UUID, error)
	DeleteSession(session_id string, ctx context.Context) error
	Login(username string, password string, ctx context.Context) (*string, error)
	Register(username string, password string, ctx context.Context) error
}
