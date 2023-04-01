package controllers

import (
	"bytes"
	"encoding/json"
	"phota/internal/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type postServiceHandler struct {
	PostService      entities.PostService
	FollowingService entities.FollowingService
}

type decoding struct {
	Visibility  string
	Description string
}

func NewPostServiceHandler(
	app *fiber.App,
	p entities.PostService,
	f entities.FollowingService,
) {
	handler := &postServiceHandler{p, f}
	app.Post("/post", handler.PostHandler)
	app.Get("/post/:post_id<guid>", handler.GettingPostHandler)
	app.Patch("/post/:post_id<guid>", handler.PostChangingHandler)
	app.Put("/post/:post_id<guid>/like", handler.LikeHandler)
	app.Get("/post/:post_id<guid>/likes", handler.GetLikesHandler)
	app.Delete("/post/:post_id<guid>/like", handler.UnlikeHandler)
	app.Delete("/post/:post_id<guid>", handler.DeletePostHandler)
}

func (p *postServiceHandler) PostHandler(c *fiber.Ctx) error {
	if _, ok := c.Locals("user_id").(*uuid.UUID); !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	description := c.FormValue("description")
	user_id, _ := c.Locals("user_id").(*uuid.UUID)
	post_id, err := p.PostService.Post(*user_id, c.Context(), description)
	if err != nil {
		return error_handling(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"post_id": post_id})
}

func (p *postServiceHandler) GettingPostHandler(c *fiber.Ctx) error {
	if _, ok := c.Locals("user_id").(*uuid.UUID); !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	post_id_raw := c.Params("post_id")
	user_id, _ := c.Locals("user_id").(*uuid.UUID)
	post_id, err := uuid.Parse(post_id_raw)
	if err != nil {
		return err
	}
	post, err := p.PostService.GetPost(post_id, *user_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"description": post.Description,
		"created at":  post.CreatedAt,
		"post_id":     post.PostID,
		"visibility":  post.Visibility,
	})
}

func (p *postServiceHandler) PostChangingHandler(c *fiber.Ctx) error {
	if _, ok := c.Locals("user_id").(*uuid.UUID); !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	decoder := json.NewDecoder(bytes.NewReader(c.Body()))
	var data decoding
	err := decoder.Decode(&data)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	visibility := data.Visibility
	description := data.Description
	user_id, _ := c.Locals("user_id").(*uuid.UUID)
	post_id_raw := c.Params("post_id")
	post_id, err := uuid.Parse(post_id_raw)
	if err != nil {
		return err
	}
	err = p.PostService.ChangePost(entities.ChangePostParams{
		Visibility:  visibility,
		Description: description,
		PostID:      post_id,
		UserID:      *user_id,
	}, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (p *postServiceHandler) LikeHandler(c *fiber.Ctx) error {
	if _, ok := c.Locals("user_id").(*uuid.UUID); !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	post_id_raw := c.Params("post_id")
	user_id, _ := c.Locals("user_id").(*uuid.UUID)
	post_id, err := uuid.Parse(post_id_raw)
	if err != nil {
		return err
	}
	err = p.PostService.Like(post_id, *user_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (p *postServiceHandler) GetLikesHandler(c *fiber.Ctx) error {
	if _, ok := c.Locals("user_id").(*uuid.UUID); !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	post_id_raw := c.Params("post_id")
	user_id, _ := c.Locals("user_id").(*uuid.UUID)
	post_id, err := uuid.Parse(post_id_raw)
	if err != nil {
		return err
	}
	likes, err := p.PostService.GetLikes(post_id, *user_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"likes": likes})
}

func (p *postServiceHandler) UnlikeHandler(c *fiber.Ctx) error {
	if _, ok := c.Locals("user_id").(*uuid.UUID); !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	post_id_raw := c.Params("post_id")
	user_id, _ := c.Locals("user_id").(*uuid.UUID)
	post_id, err := uuid.Parse(post_id_raw)
	if err != nil {
		return err
	}
	err = p.PostService.Unlike(post_id, *user_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (p *postServiceHandler) DeletePostHandler(c *fiber.Ctx) error {
	if _, ok := c.Locals("user_id").(*uuid.UUID); !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	user_id, _ := c.Locals("user_id").(*uuid.UUID)
	post_id_raw := c.Params("post_id")
	post_id, err := uuid.Parse(post_id_raw)
	if err != nil {
		return err
	}
	err = p.PostService.DeletePost(post_id, *user_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}
