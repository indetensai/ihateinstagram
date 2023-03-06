package usecases

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"mime/multipart"
	"phota/internal/entities"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type imageService struct {
	repo entities.ImageRepository
}

func NewImageService(repo entities.ImageRepository) entities.ImageService {
	return &imageService{repo: repo}
}

func (i *imageService) UploadImage(
	ctx context.Context,
	post_id uuid.UUID,
	user_id uuid.UUID,
	content *multipart.FileHeader,
) error {
	pic, err := content.Open()
	if err != nil {
		return err
	}
	img, _, err := image.Decode(pic)
	if err != nil {
		return err
	}
	eg, _ := errgroup.WithContext(ctx)
	img_buf := new(bytes.Buffer)
	thumbnail_buf := new(bytes.Buffer)
	eg.Go(func() error {
		return jpeg.Encode(img_buf, img, nil)
	})
	eg.Go(func() error {
		thumbnail := imaging.Resize(img, 128, 128, imaging.Lanczos)
		return jpeg.Encode(thumbnail_buf, thumbnail, nil)
	})
	if eg.Wait() != nil {
		return err
	}
	err = i.repo.UploadImage(ctx, post_id, user_id, img_buf.Bytes(), thumbnail_buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (i *imageService) GetImages(ctx context.Context, post_id uuid.UUID) (*[][]byte, error) {
	images, err := i.repo.GetImages(ctx, post_id)
	if err != nil {
		return nil, err
	}
	return &images, nil
}

func (i *imageService) GetThumbnails(ctx context.Context, post_id uuid.UUID) (*[][]byte, error) {
	thumbnails, err := i.repo.GetThumbnails(ctx, post_id)
	if err != nil {
		return nil, err
	}
	return &thumbnails, nil
}
