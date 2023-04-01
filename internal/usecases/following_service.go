package usecases

import (
	"context"
	"phota/internal/entities"

	"github.com/google/uuid"
)

type followingService struct {
	followingRepo entities.FollowingRepository
	userService   entities.UserService
}

func NewFollowingService(f entities.FollowingRepository, u entities.UserService) entities.FollowingService {
	return &followingService{f, u}
}

func (f *followingService) Follow(
	follower_id uuid.UUID,
	user_id uuid.UUID,
	ctx context.Context) error {
	err := f.followingRepo.CreateFollowing(ctx, follower_id, user_id)
	if err != nil {
		return err
	}
	return nil
}
func (f *followingService) Unfollow(
	follower_id uuid.UUID,
	user_id uuid.UUID,
	ctx context.Context) error {
	err := f.followingRepo.DeleteFollowing(ctx, follower_id, user_id)
	if err != nil {
		return err
	}
	return nil
}

func (f *followingService) GetFollowers(
	user_id uuid.UUID,
	ctx context.Context,
) ([]entities.AppUser, error) {
	followers, err := f.followingRepo.GetFollowers(ctx, user_id)
	if err != nil {
		return nil, err
	}
	return followers, nil
}

func (f *followingService) IsFollowing(
	user_id uuid.UUID,
	follower_id uuid.UUID,
	ctx context.Context,
) bool {
	check := f.followingRepo.IsFollowing(user_id, follower_id, ctx)
	return check
}
