package entities

import (
	"context"
	"time"

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
}
