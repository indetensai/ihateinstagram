package usecases

import (
	"context"
	"phota/internal/entities"
	"phota/internal/usecases/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type followingService struct {
	db *pgx.Conn
}

func (f *followingService) Follow(
	follower_id uuid.UUID,
	user_id uuid.UUID,
	session_id string,
	ctx context.Context) error {
	var user_service userService
	session_follower_id, err := user_service.CheckSession(session_id, ctx)
	if err != nil {
		return err
	}
	if follower_id != *session_follower_id {
		return entities.ErrNotAuthorized
	}
	following_repository := repository.NewFollowingRepository(f.db)
	err = following_repository.Follow(ctx, follower_id, user_id)
	if err != nil {
		return err
	}
	return nil
}
func (f *followingService) Unfollow(
	follower_id uuid.UUID,
	user_id uuid.UUID,
	session_id string,
	ctx context.Context) error {
	var user_service userService
	session_follower_id, err := user_service.CheckSession(session_id, ctx)
	if err != nil {
		return err
	}
	if follower_id != *session_follower_id {
		return entities.ErrNotAuthorized
	}
	following_repository := repository.NewFollowingRepository(f.db)
	err = following_repository.Unfollow(ctx, follower_id, user_id)
	if err != nil {
		return err
	}
	return nil
}
func (f *followingService) GetFollowers(
	user_id uuid.UUID,
	ctx context.Context,
) (*[]uuid.UUID, error) {
	following_repository := repository.NewFollowingRepository(f.db)
	followers, err := following_repository.GetFollowers(ctx, user_id)
	if err != nil {
		return nil, err
	}
	return &followers, nil
}
