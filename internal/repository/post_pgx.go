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

func (p *PostRepository) CreatePost(user_id uuid.UUID, ctx context.Context, desription string) (*uuid.UUID, error) {
	var post_id uuid.UUID
	err := p.db.QueryRow(
		ctx,
		"INSERT INTO post (user_id,description) VALUES ($1,$2) RETURNING post_id",
		user_id,
		desription,
	).Scan(&post_id)
	if err != nil {
		return nil, entities.ErrDuplicate
	}
	return &post_id, nil
}

func (p *PostRepository) GetPostByID(
	user_id uuid.UUID,
	post_id uuid.UUID,
	ctx context.Context,
) (*entities.Post, error) {
	var post entities.Post
	err := p.db.QueryRow(
		ctx,
		"SELECT * FROM post WHERE post_id=$1",
		post_id,
	).Scan(
		&post.PostID, &post.UserID, &post.Description, &post.Visibility, &post.CreatedAt,
	)
	if err != nil {
		return nil, entities.ErrNotFound
	}
	return &post, nil
}

func (p *PostRepository) ChangePost(content entities.ChangePostParams, ctx context.Context) error {
	_, err := p.db.Exec(
		ctx,
		"UPDATE post SET vision=$1,description=$2 WHERE post_id=$3",
		content.Visibility,
		content.Description,
		content.PostID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostRepository) CreateLike(user_id uuid.UUID, post_id uuid.UUID, ctx context.Context) error {
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
		"INSERT INTO likes (user_id,post_id) VALUES ($1,$2)",
		user_id,
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
	post_id uuid.UUID,
	ctx context.Context,
) ([]entities.AppUser, error) {
	rows, err := p.db.Query(
		ctx,
		`SELECT app_user.user_id, app_user.username 
		FROM app_user
		INNER JOIN likes ON app_user.user_id=like.user_id
		WHERE like.post_id=$1`,
		post_id,
	)
	if err != nil {
		return nil, entities.ErrNotFound
	}
	defer rows.Close()
	var likes []entities.AppUser
	for rows.Next() {
		var u entities.AppUser
		err = rows.Scan(&u)
		if err != nil {
			return nil, err
		}
		likes = append(likes, u)
	}
	return likes, nil
}

func (p *PostRepository) DeleteLike(user_id uuid.UUID, post_id uuid.UUID, ctx context.Context) error {
	_, err := p.db.Exec(
		ctx,
		"DELETE FROM likes WHERE user_id=$1 AND post_id=$2",
		user_id,
		post_id,
	)
	return err
}

func (p *PostRepository) DeletePost(post_id uuid.UUID, ctx context.Context) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = p.db.Exec(
		ctx,
		"DELETE FROM image WHERE post_id=$1",
		post_id,
	)
	if err != nil {
		return err
	}
	_, err = p.db.Exec(
		ctx,
		"DELETE FROM post WHERE post_id=$1",
		post_id,
	)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	return err
}
