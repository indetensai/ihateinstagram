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

type Post struct {
	PostID      uuid.UUID
	UserID      uuid.UUID
	Description string
	Visibility  string
	CreatedAt   time.Time
}

type Image struct {
	ImageID uuid.UUID
	UserID  uuid.UUID
	PostID  uuid.UUID
	Picture []byte
}

type Following struct {
	FollowerId     uuid.UUID
	UserID         uuid.UUID
	FollowingSince time.Time
}

type Like struct {
	UserID  uuid.UUID
	PostID  uuid.UUID
	LikedAt time.Time
}

type UserService interface {
	CheckSession(session_id string, ctx context.Context) (*uuid.UUID, error)
	DeleteSession(session_id string, ctx context.Context) error
	Login(username string, password string, ctx context.Context) error
	Register(username string, password string, ctx context.Context) error
}

type PostingService interface {
	Post(session_id string, ctx context.Context, desription string) error
	Image(*fiber.Ctx) error
	GettingPost(*fiber.Ctx) error
	PostChanging(*fiber.Ctx) error
	Like(*fiber.Ctx) error
	GetLikes(*fiber.Ctx) error
	Unlike(*fiber.Ctx) error
}

type FollowingService interface {
	Follow(*fiber.Ctx) error
	Unfollow(*fiber.Ctx) error
	GetFollowers(*fiber.Ctx) error
}
