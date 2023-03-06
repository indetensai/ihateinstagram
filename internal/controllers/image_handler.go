package controllers

import (
	"bytes"
	"image"
	"image/jpeg"
	"phota/internal/entities"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type imageServiceHandler struct {
	ImageService entities.ImageService
	PostService  entities.PostService
}

func NewImageServiceHandler(app *fiber.App, i entities.ImageService, p entities.PostService) {
	handler := &imageServiceHandler{i, p}
	app.Post("/post/:post_id<guid>/image", handler.UploadImageHandler)
}

func (i *imageServiceHandler) UploadImageHandler(c *fiber.Ctx) error {
	post_id, err := uuid.Parse(c.Params("post_id"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	user_id, _ := c.Locals("user_id").(*uuid.UUID)
	if user_id == nil {
		return entities.ErrNotAuthorized
	}
	post, err := i.PostService.GettingPost(post_id, user_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	if post.UserID != *user_id {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	content, err := c.FormFile("image")
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	pic, err := content.Open()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	img, _, err := image.Decode(pic)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	eg, _ := errgroup.WithContext(c.Context())
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
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	err = i.ImageService.UploadImage(c.Context(), post_id, *user_id, content)
	if err != nil {
		return error_handling(c, err)
	}
	return nil
}
