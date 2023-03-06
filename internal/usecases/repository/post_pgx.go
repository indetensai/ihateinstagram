package repository

import (
	"context"
	"phota/internal/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostRepository struct {
	db *pgx.Conn
}

func NewPostRepository(db *pgx.Conn) entities.PostRepository {
	return &PostRepository{db: db}
}

func (p *PostRepository) Post(user_id uuid.UUID, ctx context.Context, desription string) (*uuid.UUID, error) {
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
		return nil, entities.ErrDuplicate
	}
	var post_id uuid.UUID
	err = tx.QueryRow(
		ctx,
		"SELECT post_id FROM post ORDER BY created_at DESC",
	).Scan(&post_id)
	if err != nil {
		return nil, entities.ErrNotFound
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return &post_id, nil
}

func (p *PostRepository) GettingPost(user_id *uuid.UUID, post_id uuid.UUID, ctx context.Context) (*entities.Post, error) {
	var post entities.Post
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
		&post.PostID, &post.UserID, &post.Description, &post.Visibility, &post.CreatedAt,
	)
	if err != nil {
		return nil, entities.ErrNotFound
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (p *PostRepository) PostChanging(
	user_id *uuid.UUID,
	visibility string,
	description string,
	post_id uuid.UUID,
	ctx context.Context) error {
	var received_post entities.Post
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
		&received_post.PostID,
		&received_post.UserID,
		&received_post.Description,
		&received_post.Visibility,
		&received_post.CreatedAt,
	)
	if err != nil {
		return err
	}
	if *user_id != received_post.UserID {
		return entities.ErrNotAuthorized
	}
	check := received_post
	if description != check.Description && description != "" {
		check.Description = description
	}
	if check.Visibility != visibility && visibility != "" {
		check.Visibility = visibility
	}
	if check != received_post {
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
	if err != nil {
		return err
	}
	return nil
}

func (p *PostRepository) Like(user_id *uuid.UUID, post_id uuid.UUID, ctx context.Context) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(
		ctx,
		"SELECT 1 FROM post WHERE post_id=$1",
		post_id,
	)
	if err != nil {
		return entities.ErrNotFound
	}
	_, err = tx.Exec(
		ctx,
		"INSERT INTO liking (user_id,post_id) VALUES ($1,$2)",
		*user_id,
		post_id,
	)
	if err != nil {
		return entities.ErrDuplicate
	}
	tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostRepository) GetLikes(
	user_id *uuid.UUID,
	post_id uuid.UUID,
	ctx context.Context,
) (*[]uuid.UUID, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	rows, err := tx.Query(
		ctx,
		"SELECT user_id FROM liking WHERE post_id=$1",
		post_id,
	)
	if err != nil {
		return nil, entities.ErrNotFound
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

func (p *PostRepository) Unlike(user_id *uuid.UUID, post_id uuid.UUID, ctx context.Context) error {
	_, err := p.db.Exec(
		ctx,
		"DELETE FROM liking WHERE user_id=$1 AND post_id=$2",
		*user_id,
		post_id,
	)
	return err
}

func (p *PostRepository) DeletePost(post_id uuid.UUID, ctx context.Context) error {
	_, err := p.db.Exec(
		ctx,
		"DELETE FROM post WHERE post_id=$1",
		post_id,
	)
	return err
}
