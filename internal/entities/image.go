package entities

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
)

type Image struct {
	ImageID uuid.UUID
	UserID  uuid.UUID
	PostID  uuid.UUID
	Content []byte
}

type ImageService interface {
	UploadImage(ctx context.Context, post_id uuid.UUID, user_id uuid.UUID, content *multipart.FileHeader) error
	GetImages(ctx context.Context, post_id uuid.UUID) (*[][]byte, error)
	GetThumbnails(ctx context.Context, post_id uuid.UUID) (*[][]byte, error)
}
type ImageRepository interface {
	UploadImage(ctx context.Context, post_id uuid.UUID, user_id uuid.UUID, content []byte, thumbnail []byte) error
	GetImages(ctx context.Context, post_id uuid.UUID) ([][]byte, error)
	GetThumbnails(ctx context.Context, post_id uuid.UUID) ([][]byte, error)
}
