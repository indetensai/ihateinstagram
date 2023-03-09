package entities

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Post struct {
	PostID      uuid.UUID
	UserID      uuid.UUID
	Description string
	Visibility  string
	CreatedAt   time.Time
}

type Like struct {
	UserID  uuid.UUID
	PostID  uuid.UUID
	LikedAt time.Time
}

type PostHandler interface {
	PostHandler(c *fiber.Ctx) error
	GettingPostHandler(c *fiber.Ctx) error
	PostChangingHandler(c *fiber.Ctx) error
	LikeHandler(c *fiber.Ctx) error
	GetLikesHandler(c *fiber.Ctx) error
	UnlikeHandler(c *fiber.Ctx) error
}

type PostService interface {
	Post(user_id *uuid.UUID, ctx context.Context, desription string) (*uuid.UUID, error)
	GettingPost(post_id uuid.UUID, user_id *uuid.UUID, ctx context.Context) (*Post, error)
	PostChanging(
		visibility string,
		description string,
		post_id uuid.UUID,
		user_id *uuid.UUID,
		ctx context.Context) error
	Like(post_id uuid.UUID, user_id *uuid.UUID, ctx context.Context) error
	GetLikes(user_id *uuid.UUID, post_id uuid.UUID, ctx context.Context) (*[]uuid.UUID, error)
	Unlike(user_id *uuid.UUID, post_id uuid.UUID, ctx context.Context) error
	DeletePost(post_id uuid.UUID, ctx context.Context) error
}

type PostRepository interface {
	Post(user_id uuid.UUID, ctx context.Context, desription string) (*uuid.UUID, error)
	GettingPost(user_id *uuid.UUID, post_id uuid.UUID, ctx context.Context) (*Post, error)
	PostChanging(
		user_id *uuid.UUID,
		visibility string,
		description string,
		post_id uuid.UUID,
		ctx context.Context) error
	Like(user_id *uuid.UUID, post_id uuid.UUID, ctx context.Context) error
	GetLikes(user_id *uuid.UUID, post_id uuid.UUID, ctx context.Context) (*[]uuid.UUID, error)
	Unlike(user_id *uuid.UUID, post_id uuid.UUID, ctx context.Context) error
	DeletePost(post_id uuid.UUID, ctx context.Context) error
}
