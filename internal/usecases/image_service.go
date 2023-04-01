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
	imageRepo        entities.ImageRepository
	postService      entities.PostService
	followingService entities.FollowingService
}

func NewImageService(
	repo entities.ImageRepository,
	postService entities.PostService,
	followingService entities.FollowingService,
) entities.ImageService {
	return &imageService{imageRepo: repo, postService: postService, followingService: followingService}
}

func (i *imageService) UploadImage(
	ctx context.Context,
	post_id uuid.UUID,
	user_id uuid.UUID,
	content *multipart.FileHeader,
) error {
	post, err := i.postService.GetPost(post_id, user_id, context.Background())
	if err != nil {
		return err
	}
	if post.UserID != user_id {
		return entities.ErrNotAuthorized
	}
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
	err = i.imageRepo.CreateImage(ctx, post_id, user_id, img_buf.Bytes(), thumbnail_buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (i *imageService) GetImages(ctx context.Context, post_id uuid.UUID, user_id uuid.UUID) ([][]byte, error) {
	post, err := i.postService.GetPost(post_id, user_id, context.Background())
	if err != nil {
		return nil, err
	}
	switch post.Visibility {
	case "followers":
		if post.UserID != user_id || !i.followingService.IsFollowing(post.UserID, user_id, context.Background()) {
			return nil, entities.ErrNotAuthorized
		}
	case "private":
		if post.UserID != user_id {
			return nil, entities.ErrNotAuthorized
		}
	}
	images, err := i.imageRepo.GetImages(ctx, post_id)
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (i *imageService) GetThumbnails(ctx context.Context, post_id uuid.UUID, user_id uuid.UUID) ([][]byte, error) {
	post, err := i.postService.GetPost(post_id, user_id, context.Background())
	if err != nil {
		return nil, err
	}
	switch post.Visibility {
	case "followers":
		if post.UserID != user_id || !i.followingService.IsFollowing(post.UserID, user_id, context.Background()) {
			return nil, entities.ErrNotAuthorized
		}
	case "private":
		if post.UserID != user_id {
			return nil, entities.ErrNotAuthorized
		}
	}
	thumbnails, err := i.imageRepo.GetThumbnails(ctx, post_id)
	if err != nil {
		return nil, err
	}
	return thumbnails, nil
}
