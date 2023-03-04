package controllers

import (
	"phota/internal/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PostServiceHandler struct {
	PostService entities.PostService
}

func NewPostServiceHandler(app *fiber.App, p entities.PostService) {
	handler := PostServiceHandler{p}
	app.Post("/post", handler.PostHandler)
	app.Get("/post/:post_id<guid>", handler.GettingPostHandler)
	app.Patch("/post/:post_id<guid>", handler.PostChangingHandler)
	app.Put("/post/:post_id<guid>/like", handler.LikeHandler)
	app.Get("/post/:post_id<guid>/likes", handler.GetLikesHandler)
	app.Delete("/post/:post_id<guid>/like", handler.UnlikeHandler)
}

func (p *PostServiceHandler) PostHandler(c *fiber.Ctx) error {
	session_id := c.FormValue("session_id")
	description := c.FormValue("description")
	post_id, err := p.PostService.Post(session_id, c.Context(), description)
	if err != nil {
		return error_handling(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"post_id": post_id})
}

func (p *PostServiceHandler) GettingPostHandler(c *fiber.Ctx) error {
	post_id_raw := c.Params("post_id")
	session_id := c.FormValue("session_id")
	post_id, err := uuid.Parse(post_id_raw)
	if err != nil {
		return err
	}
	post, err := p.PostService.GettingPost(post_id, session_id, c.Context())
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

func (p *PostServiceHandler) PostChangingHandler(c *fiber.Ctx) error {
	visibility := c.FormValue("visibility")
	description := c.FormValue("description")
	post_id_raw := c.Params("post_id")
	session_id := c.FormValue("session_id")
	post_id, err := uuid.Parse(post_id_raw)
	if err != nil {
		return err
	}
	err = p.PostService.PostChanging(visibility, description, post_id, session_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (p *PostServiceHandler) LikeHandler(c *fiber.Ctx) error {
	post_id_raw := c.Params("post_id")
	session_id := c.FormValue("session_id")
	post_id, err := uuid.Parse(post_id_raw)
	if err != nil {
		return err
	}
	err = p.PostService.Like(post_id, session_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (p *PostServiceHandler) GetLikesHandler(c *fiber.Ctx) error {
	post_id_raw := c.Params("post_id")
	session_id := c.FormValue("session_id")
	post_id, err := uuid.Parse(post_id_raw)
	if err != nil {
		return err
	}
	likes, err := p.PostService.GetLikes(session_id, post_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"likes": likes})
}

func (p *PostServiceHandler) UnlikeHandler(c *fiber.Ctx) error {
	post_id_raw := c.Params("post_id")
	session_id := c.FormValue("session_id")
	post_id, err := uuid.Parse(post_id_raw)
	if err != nil {
		return err
	}
	err = p.PostService.Unlike(session_id, post_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}
