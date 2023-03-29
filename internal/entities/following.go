package entities

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Following struct {
	UserID         uuid.UUID
	FollowingSince time.Time
}

type FollowingHandler interface {
	FollowHandler(c *fiber.Ctx) error
	UnfollowHandler(c *fiber.Ctx) error
	GetFollowersHandler(c *fiber.Ctx) error
}

type FollowingService interface {
	Follow(follower_id uuid.UUID, user_id uuid.UUID, ctx context.Context) error
	Unfollow(follower_id uuid.UUID, user_id uuid.UUID, ctx context.Context) error
	GetFollowers(user_id uuid.UUID, ctx context.Context) (*[]uuid.UUID, error)
	IsFollowing(user_id uuid.UUID, follower_id uuid.UUID, ctx context.Context) bool
}

type FollowingRepository interface {
	Follow(ctx context.Context, follower_id uuid.UUID, user_id uuid.UUID) error
	Unfollow(ctx context.Context, follower_id uuid.UUID, user_id uuid.UUID) error
	GetFollowers(ctx context.Context, user_id uuid.UUID) ([]uuid.UUID, error)
	IsFollowing(user_id uuid.UUID, follower_id uuid.UUID, ctx context.Context) bool
}
