package usecases

import (
	"context"
	"phota/internal/entities"

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
	_, err = f.db.Exec(
		ctx,
		"INSERT INTO following (follower_id,user_id) VALUES ($1,$2)",
		follower_id,
		user_id,
	)
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
	_, err = f.db.Exec(
		ctx,
		"DELETE FROM following WHERE user_id=$1 AND follower_id=$2",
		user_id,
		follower_id,
	)
	if err != nil {
		return err
	}
	return nil
}
func (f *followingService) GetFollowers(
	user_id uuid.UUID,
	ctx context.Context,
) (*[]uuid.UUID, error) {

	rows, err := f.db.Query(
		ctx,
		"SELECT follower_id FROM following WHERE user_id=$1",
		user_id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var followers []uuid.UUID
	for rows.Next() {
		var r uuid.UUID
		err = rows.Scan(&r)
		if err != nil {
			return nil, err
		}
		followers = append(followers, r)
	}
	return &followers, nil
}
