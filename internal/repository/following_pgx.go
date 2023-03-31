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

func (f *FollowingRepository) CreateFollowing(
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
		return entities.ErrDuplicate
	}
	return nil
}

func (f *FollowingRepository) DeleteFollowing(
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

func (f *FollowingRepository) GetFollowers(ctx context.Context, user_id uuid.UUID) ([]entities.AppUser, error) {
	rows, err := f.db.Query(
		ctx,
		`SELECT app_user.user_id, app_user.username 
		FROM app_user
		INNER JOIN following ON app_user.user_id=following.follower_id
		WHERE following.user_id=$1`,
		user_id,
	)
	if err != nil {
		return nil, entities.ErrNotFound
	}
	defer rows.Close()
	var followers []entities.AppUser
	for rows.Next() {
		var r entities.AppUser
		err = rows.Scan(&r.UserID, &r.Username)
		if err != nil {
			return nil, err
		}
		followers = append(followers, r)
	}
	return followers, nil
}

func (f *FollowingRepository) IsFollowing(
	user_id uuid.UUID,
	follower_id uuid.UUID,
	ctx context.Context,
) bool {
	_, err := f.db.Exec(
		ctx,
		"SELECT EXISTS (SELECT 1 FROM following WHERE user_id=$1 AND follower_id=$2)",
		user_id,
		follower_id,
	)
	if err != nil {
		return false
	}
	return true
}
