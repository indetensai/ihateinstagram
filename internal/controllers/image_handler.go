package controllers

import (
	"phota/internal/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type imageServiceHandler struct {
	ImageService entities.ImageService
}

func NewImageServiceHandler(
	app *fiber.App,
	i entities.ImageService,
) {
	handler := &imageServiceHandler{i}
	app.Post("/post/:post_id<guid>/image", handler.UploadImageHandler)
	app.Get("/post/:post_id<guid>/images", handler.GetImages)
	app.Get("/post/:post_id<guid>/thumbnails", handler.GetThumbnails)
}

func (i *imageServiceHandler) UploadImageHandler(c *fiber.Ctx) error {
	user_id, ok := c.Locals("user_id").(*uuid.UUID)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	post_id, err := uuid.Parse(c.Params("post_id"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	content, err := c.FormFile("image")
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	err = i.ImageService.UploadImage(c.Context(), post_id, *user_id, content)
	if err != nil {
		return error_handling(c, err)
	}
	return nil
}

func (i *imageServiceHandler) GetImages(c *fiber.Ctx) error {
	user_id, ok := c.Locals("user_id").(*uuid.UUID)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	post_id, err := uuid.Parse(c.Params("post_id"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	images, err := i.ImageService.GetImages(c.Context(), post_id, *user_id)
	if err != nil {
		return error_handling(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"images": images})
}

func (i *imageServiceHandler) GetThumbnails(c *fiber.Ctx) error {
	user_id, ok := c.Locals("user_id").(*uuid.UUID)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	post_id, err := uuid.Parse(c.Params("post_id"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	thumbnails, err := i.ImageService.GetThumbnails(c.Context(), post_id, *user_id)
	if err != nil {
		return error_handling(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"thumbnails": thumbnails})
}
