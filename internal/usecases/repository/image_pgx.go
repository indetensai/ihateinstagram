package repository

import (
	"context"
	"phota/internal/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type imageRepository struct {
	db *pgx.Conn
}

func NewImageRepository(db *pgx.Conn) entities.ImageRepository {
	return &imageRepository{db: db}
}

func (i *imageRepository) UploadImage(
	ctx context.Context,
	post_id uuid.UUID,
	user_id uuid.UUID,
	content []byte,
	thumbnail []byte,
) error {
	_, err := i.db.Exec(
		ctx,
		"INSERT INTO image (post_id,user_id,content,thumbnail) VALUES ($1,$2,$3,$4)",
		post_id,
		user_id,
		content,
		thumbnail,
	)
	if err != nil {
		return err
	}
	return nil
}
