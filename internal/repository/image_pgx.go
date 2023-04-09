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

func (i *imageRepository) CreateImage(
	ctx context.Context,
	post_id uuid.UUID,
	user_id uuid.UUID,
	content []byte,
	thumbnail []byte,
) error {
	_, err := i.db.Exec(
		ctx,
		"INSERT INTO images (post_id,user_id,content,thumbnail) VALUES ($1,$2,$3,$4)",
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

func (i *imageRepository) GetImages(
	ctx context.Context,
	post_id uuid.UUID,
) ([][]byte, error) {
	rows, err := i.db.Query(
		ctx,
		"SELECT content FROM images WHERE post_id=$1",
		post_id,
	)
	if err != nil {
		return nil, entities.ErrNotFound
	}
	defer rows.Close()
	var images [][]byte
	for rows.Next() {
		var i []byte
		err = rows.Scan(&i)
		if err != nil {
			return nil, err
		}
		images = append(images, i)
	}
	return images, nil
}

func (i *imageRepository) GetThumbnails(
	ctx context.Context,
	post_id uuid.UUID,
) ([][]byte, error) {
	rows, err := i.db.Query(
		ctx,
		"SELECT thumbnail FROM images WHERE post_id=$1",
		post_id,
	)
	if err != nil {
		return nil, entities.ErrNotFound
	}
	defer rows.Close()
	var thumbnails [][]byte
	for rows.Next() {
		var i []byte
		err = rows.Scan(&i)
		if err != nil {
			return nil, err
		}
		thumbnails = append(thumbnails, i)
	}
	return thumbnails, nil
}
