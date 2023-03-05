package controllers

import (
	"bytes"
	"encoding/json"
	"phota/internal/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PostServiceHandler struct {
	PostService      entities.PostService
	FollowingService entities.FollowingService
}

type decoding struct {
	Visibility  string
	Description string
}

func NewPostServiceHandler(app *fiber.App, p entities.PostService, f entities.FollowingService) {
	handler := PostServiceHandler{p, f}
	app.Post("/post", handler.PostHandler)
	app.Get("/post/:post_id<guid>", handler.GettingPostHandler)
	app.Patch("/post/:post_id<guid>", handler.PostChangingHandler)
	app.Put("/post/:post_id<guid>/like", handler.LikeHandler)
	app.Get("/post/:post_id<guid>/likes", handler.GetLikesHandler)
	app.Delete("/post/:post_id<guid>/like", handler.UnlikeHandler)
}

func (p *PostServiceHandler) PostHandler(c *fiber.Ctx) error {
	description := c.FormValue("description")
	user_id, _ := c.Locals("user_id").(*uuid.UUID)
	if user_id == nil {
		return entities.ErrNotAuthorized
	}
	post_id, err := p.PostService.Post(user_id, c.Context(), description)
	if err != nil {
		return error_handling(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"post_id": post_id})
}

func (p *PostServiceHandler) GettingPostHandler(c *fiber.Ctx) error {
	post_id_raw := c.Params("post_id")
	user_id, _ := c.Locals("user_id").(*uuid.UUID)
	if user_id == nil {
		return entities.ErrNotAuthorized
	}
	post_id, err := uuid.Parse(post_id_raw)
	if err != nil {
		return err
	}
	post, err := p.PostService.GettingPost(post_id, user_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	switch post.Visibility {
	case "followers":
		if !p.FollowingService.IsFollowing(post.UserID, *user_id, c.Context()) || post.UserID != *user_id {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	case "private":
		if post.UserID != *user_id {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"description": post.Description,
		"created at":  post.CreatedAt,
		"post_id":     post.PostID,
		"visibility":  post.Visibility,
	})
}

func (p *PostServiceHandler) PostChangingHandler(c *fiber.Ctx) error {
	decoder := json.NewDecoder(bytes.NewReader(c.Body()))
	var result decoding
	err := decoder.Decode(&result)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	visibility := result.Visibility
	description := result.Description
	post_id_raw := c.Params("post_id")
	user_id, _ := c.Locals("user_id").(*uuid.UUID)
	if user_id == nil {
		return entities.ErrNotAuthorized
	}
	post_id, err := uuid.Parse(post_id_raw)
	if err != nil {
		return err
	}
	post, err := p.PostService.GettingPost(post_id, user_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	if post.UserID != *user_id {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	err = p.PostService.PostChanging(visibility, description, post_id, user_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (p *PostServiceHandler) LikeHandler(c *fiber.Ctx) error {
	post_id_raw := c.Params("post_id")
	user_id, _ := c.Locals("user_id").(*uuid.UUID)
	if user_id == nil {
		return entities.ErrNotAuthorized
	}
	post_id, err := uuid.Parse(post_id_raw)
	if err != nil {
		return err
	}
	post, err := p.PostService.GettingPost(post_id, user_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	switch post.Visibility {
	case "followers":
		if !p.FollowingService.IsFollowing(post.UserID, *user_id, c.Context()) || post.UserID != *user_id {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	case "private":
		if post.UserID != *user_id {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}
	err = p.PostService.Like(post_id, user_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (p *PostServiceHandler) GetLikesHandler(c *fiber.Ctx) error {
	post_id_raw := c.Params("post_id")
	user_id, _ := c.Locals("user_id").(*uuid.UUID)
	if user_id == nil {
		return entities.ErrNotAuthorized
	}
	post_id, err := uuid.Parse(post_id_raw)
	if err != nil {
		return err
	}
	post, err := p.PostService.GettingPost(post_id, user_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	switch post.Visibility {
	case "followers":
		if !p.FollowingService.IsFollowing(post.UserID, *user_id, c.Context()) || post.UserID != *user_id {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	case "private":
		if post.UserID != *user_id {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}
	likes, err := p.PostService.GetLikes(user_id, post_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"likes": likes})
}

func (p *PostServiceHandler) UnlikeHandler(c *fiber.Ctx) error {
	post_id_raw := c.Params("post_id")
	user_id, _ := c.Locals("user_id").(*uuid.UUID)
	if user_id == nil {
		return entities.ErrNotAuthorized
	}
	post_id, err := uuid.Parse(post_id_raw)
	if err != nil {
		return err
	}
	post, err := p.PostService.GettingPost(post_id, user_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	switch post.Visibility {
	case "followers":
		if !p.FollowingService.IsFollowing(post.UserID, *user_id, c.Context()) || post.UserID != *user_id {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	case "private":
		if post.UserID != *user_id {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}
	err = p.PostService.Unlike(user_id, post_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}
