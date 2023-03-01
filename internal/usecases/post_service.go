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

func (p *postService) GettingPost(post_id uuid.UUID, session_id string, ctx context.Context) (*entities.Post, error) {
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
func (p *postService) PostChanging(
	visibility string,
	description string,
	post_id uuid.UUID,
	session_id string,
	ctx context.Context,
) error {
	var user_service userService
	user_id, err := user_service.CheckSession(session_id, ctx)
	if err != nil {
		return err
	}
	received_post := new(entities.Post)
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	err = tx.QueryRow(
		ctx,
		"SELECT * FROM post WHERE post_id=$1",
		post_id,
	).Scan(
		&received_post,
	)
	if *user_id != received_post.UserID {
		return entities.ErrNotAuthorized
	}
	check := *received_post
	if description != check.Description && description != "" {
		check.Description = description
	}
	if check.Visibility != visibility && visibility != "" {
		check.Visibility = visibility
	}
	if check != *received_post {
		_, err = tx.Exec(
			ctx,
			"UPDATE post SET vision=$1,description=$2 WHERE post_id=$3",
			check.Visibility,
			check.Description,
			post_id,
		)
		if err != nil {
			return err
		}
	}
	tx.Commit(ctx)
	return nil
}
func (p *postService) Like(post_id uuid.UUID, session_id string, ctx context.Context) error {
	var user_service userService
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	user_id, err := user_service.CheckSession(session_id, ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		ctx,
		"SELECT 1 FROM post WHERE post_id=$1",
		post_id,
	)
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		ctx,
		"INSERT INTO liking (user_id,post_id) VALUES ($1,$2)",
		*user_id,
		post_id,
	)
	if err != nil {
		return err
	}
	tx.Commit(ctx)
	return nil
}
func (p *postService) GetLikes(session_id string, post_id uuid.UUID, ctx context.Context) (*[]uuid.UUID, error) {
	var user_service userService
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	user_id, err := user_service.CheckSession(session_id, ctx)
	if err != nil {
		return nil, err
	}
	var visibility string
	var owner_id uuid.UUID
	err = tx.QueryRow(
		ctx,
		"SELECT vision,user_id FROM post WHERE post_id=$1",
		post_id,
	).Scan(&visibility, &owner_id)
	if err != nil {
		return nil, err
	}
	/*if visibility != "public" && user_id == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if visibility == "private" && owner_id != *user_id {
		return c.SendStatus(fiber.StatusForbidden)
	}
	if visibility == "followers" && owner_id != *user_id {

	}
	*/
	_, err = tx.Exec(
		ctx,
		"SELECT 1 FROM following WHERE follower_id=$1 AND user_id=$2",
		*user_id,
		owner_id,
	)
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query(
		ctx,
		"SELECT user_id FROM liking WHERE post_id=$1",
		post_id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var likes []uuid.UUID
	for rows.Next() {
		var u uuid.UUID
		err = rows.Scan(&u)
		if err != nil {
			return nil, err
		}
		likes = append(likes, u)
	}
	return &likes, nil
}
func (p *postService) Unlike(session_id string, post_id uuid.UUID, ctx context.Context) error {
	var user_service userService
	user_id, err := user_service.CheckSession(session_id, ctx)
	if err != nil {
		return nil
	}
	_, err = p.db.Exec(
		ctx,
		"DELETE FROM liking WHERE user_id=$1 AND post_id=$2",
		*user_id,
		post_id,
	)
	if err != nil {
		return err
	}
	return nil
}
