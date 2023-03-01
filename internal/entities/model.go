package entities

import (
	"context"
	"time"

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
	Login(username string, password string, ctx context.Context) (*string, error)
	Register(username string, password string, ctx context.Context) error
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

type FollowingService interface {
	Follow(follower_id uuid.UUID, user_id uuid.UUID, session_id string, ctx context.Context) error
	Unfollow(follower_id uuid.UUID, user_id uuid.UUID, session_id string, ctx context.Context) error
	GetFollowers(user_id uuid.UUID, ctx context.Context) (*[]uuid.UUID, error)
}
