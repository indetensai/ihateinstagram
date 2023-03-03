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

type Image struct {
	ImageID uuid.UUID
	UserID  uuid.UUID
	PostID  uuid.UUID
	Picture []byte
}

type Like struct {
	UserID  uuid.UUID
	PostID  uuid.UUID
	LikedAt time.Time
}

type PostingService interface {
	Post(session_id string, ctx context.Context, desription string) (*uuid.UUID, error)
	GettingPost(post_id uuid.UUID, session_id string, ctx context.Context) (*Post, error)
	PostChanging(
		visibility string,
		description string,
		post_id uuid.UUID,
		session_id string,
		ctx context.Context) error
	Like(post_id uuid.UUID, session_id string, ctx context.Context) error
	GetLikes(session_id string, post_id uuid.UUID, ctx context.Context) (*[]uuid.UUID, error)
	Unlike(session_id string, post_id uuid.UUID, ctx context.Context) error
}
