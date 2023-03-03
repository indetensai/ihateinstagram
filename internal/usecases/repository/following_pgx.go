package repository

import (
	"context"
	"phota/internal/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type FollowingRepository struct {
	db *pgx.Conn
}

func NewFollowingRepository(db *pgx.Conn) entities.FollowingRepository {
	return &FollowingRepository{db: db}
}

func (f *FollowingRepository) Follow(
	ctx context.Context,
	follower_id uuid.UUID,
	user_id uuid.UUID,
) error {
	_, err := f.db.Exec(
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

func (f *FollowingRepository) Unfollow(
	ctx context.Context,
	follower_id uuid.UUID,
	user_id uuid.UUID,
) error {
	_, err := f.db.Exec(
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

func (f *FollowingRepository) GetFollowers(ctx context.Context, user_id uuid.UUID) ([]uuid.UUID, error) {
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
	return followers, nil
}
