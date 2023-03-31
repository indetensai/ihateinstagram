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

type ChangePostParams struct {
	Visibility  string
	Description string
	PostID      uuid.UUID
	UserID      uuid.UUID
}

type PostService interface {
	Post(user_id uuid.UUID, ctx context.Context, desription string) (*uuid.UUID, error)
	GetPost(post_id uuid.UUID, user_id uuid.UUID, ctx context.Context) (*Post, error)
	ChangePost(content ChangePostParams, ctx context.Context) error
	Like(post_id uuid.UUID, user_id uuid.UUID, ctx context.Context) error
	GetLikes(post_id uuid.UUID, user_id uuid.UUID, ctx context.Context) ([]AppUser, error)
	Unlike(post_id uuid.UUID, user_id uuid.UUID, ctx context.Context) error
	DeletePost(post_id uuid.UUID, user_id uuid.UUID, ctx context.Context) error
}

type PostRepository interface {
	CreatePost(user_id uuid.UUID, ctx context.Context, desription string) (*uuid.UUID, error)
	GetPostByID(user_id uuid.UUID, post_id uuid.UUID, ctx context.Context) (*Post, error)
	ChangePost(content ChangePostParams, ctx context.Context) error
	CreateLike(user_id uuid.UUID, post_id uuid.UUID, ctx context.Context) error
	GetLikes(post_id uuid.UUID, ctx context.Context) ([]AppUser, error)
	DeleteLike(user_id uuid.UUID, post_id uuid.UUID, ctx context.Context) error
	DeletePost(post_id uuid.UUID, ctx context.Context) error
}
