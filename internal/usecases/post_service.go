package usecases

import (
	"context"
	"phota/internal/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type postService struct {
	db *pgx.Conn
}

func (p *postService) Post(session_id string, ctx context.Context, desription string) (*uuid.UUID, error) {
	var user_service userService
	user_id, err := user_service.CheckSession(session_id, ctx)
	if err != nil {
		return nil, err
	}
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(
		ctx,
		"INSERT INTO post (user_id,description) VALUES ($1,$2)",
		user_id,
		desription,
	)
	if err != nil {
		return nil, err
	}
	var post_id *uuid.UUID
	err = tx.QueryRow(
		ctx,
		"SELECT post_id FROM post ORDER BY created_at DESC",
	).Scan(&post_id)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return post_id, err
}

func (p *postService) GettingPost(post_id *uuid.UUID, session_id string, ctx context.Context) (*entities.Post, error) {
	var user_service userService
	post := new(entities.Post)
	user_id, err := user_service.CheckSession(session_id, ctx)
	if err != nil {
		return nil, err
	}
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	err = tx.QueryRow(
		ctx,
		"SELECT * FROM post WHERE post_id=$1",
		post_id,
	).Scan(
		&post,
	)
	if err != nil {
		return nil, entities.ErrPostNotFound
	}
	if post.Visibility == "private" {
		if user_id == nil || *user_id != post.UserID {
			return nil, entities.ErrNotAuthorized
		}
	}
	if post.Visibility == "followers" {
		if user_id == nil {
			return nil, entities.ErrNotAuthorized
		}

		_, err = tx.Exec(
			ctx,
			"SELECT 1 FROM following WHERE follower_id=$1 AND user_id=$2",
			user_id,
			post.UserID,
		)
		if err != nil || *user_id != post.UserID {
			return nil, err
		}
	}

	tx.Commit(ctx)
	return post, nil
}
