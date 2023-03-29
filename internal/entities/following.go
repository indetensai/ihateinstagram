package entities

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Following struct {
	UserID         uuid.UUID
	FollowingSince time.Time
}
type FollowingService interface {
	Follow(follower_id uuid.UUID, user_id uuid.UUID, ctx context.Context) error
	Unfollow(follower_id uuid.UUID, user_id uuid.UUID, ctx context.Context) error
	GetFollowers(user_id uuid.UUID, ctx context.Context) ([]AppUser, error)
	IsFollowing(user_id uuid.UUID, follower_id uuid.UUID, ctx context.Context) bool
}

type FollowingRepository interface {
	CreateFollowing(ctx context.Context, follower_id uuid.UUID, user_id uuid.UUID) error
	DeleteFollowing(ctx context.Context, follower_id uuid.UUID, user_id uuid.UUID) error
	GetFollowers(ctx context.Context, user_id uuid.UUID) ([]AppUser, error)
	IsFollowing(user_id uuid.UUID, follower_id uuid.UUID, ctx context.Context) bool
}
