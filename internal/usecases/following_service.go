package usecases

import (
	"context"
	"phota/internal/entities"

	"github.com/google/uuid"
)

type followingService struct {
	repo         entities.FollowingRepository
	user_service entities.UserService
}

func NewFollowingService(f entities.FollowingRepository, u entities.UserService) entities.FollowingService {
	return &followingService{f, u}
}

func (f *followingService) Follow(
	follower_id uuid.UUID,
	user_id uuid.UUID,
	ctx context.Context) error {
	err := f.repo.Follow(ctx, follower_id, user_id)
	if err != nil {
		return err
	}
	return nil
}
func (f *followingService) Unfollow(
	follower_id uuid.UUID,
	user_id uuid.UUID,
	ctx context.Context) error {
	err := f.repo.Unfollow(ctx, follower_id, user_id)
	if err != nil {
		return err
	}
	return nil
}

func (f *followingService) GetFollowers(
	user_id uuid.UUID,
	ctx context.Context,
) (*[]uuid.UUID, error) {
	followers, err := f.repo.GetFollowers(ctx, user_id)
	if err != nil {
		return nil, err
	}
	return &followers, nil
}

func (f *followingService) IsFollowing(
	user_id uuid.UUID,
	follower_id uuid.UUID,
	ctx context.Context,
) bool {
	check := f.repo.IsFollowing(user_id, follower_id, ctx)
	return check
}
