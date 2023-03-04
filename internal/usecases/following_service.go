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
	session_id string,
	ctx context.Context) error {
	session_follower_id, err := f.user_service.CheckSession(session_id, ctx)
	if err != nil {
		return err
	}
	if follower_id != *session_follower_id {
		return entities.ErrNotAuthorized
	}
	err = f.repo.Follow(ctx, follower_id, user_id)
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
	session_follower_id, err := f.user_service.CheckSession(session_id, ctx)
	if err != nil {
		return err
	}
	if follower_id != *session_follower_id {
		return entities.ErrNotAuthorized
	}
	err = f.repo.Unfollow(ctx, follower_id, user_id)
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
