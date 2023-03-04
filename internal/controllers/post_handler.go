package controllers

import (
	"phota/internal/entities"

	"github.com/gofiber/fiber/v2"
)

type PostServiceHandler struct {
	PostService entities.PostService
}

func NewPostServiceHandler(c *fiber.App, p entities.PostService) PostServiceHandler {
	handler := &PostServiceHandler{p}
}

func (p *PostServiceHandler) PostHandler(c *fiber.Ctx) error {
	session_id := c.FormValue("session_id")
	description := c.FormValue("description")
	post_id, err := p.PostService.Post(session_id, c.Context(), description)

}
