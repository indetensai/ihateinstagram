package controllers

import (
	"phota/internal/entities"

	"github.com/gofiber/fiber/v2"
)

type UserServiceHandler struct {
	UserService entities.UserService
}

func NewUserServiceHandler(c *fiber.App, u entities.UserService) {
	handler := &UserServiceHandler{u}
	c.Post("/user/register", handler.RegisterHandler)
	//c.Post("/user/login", LoginHandler)
	//c.Delete("/session/:id<guid>", DeleteSessionHandler)
}

func (u *UserServiceHandler) RegisterHandler(c *fiber.Ctx) error {
	username := c.FormValue("username")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"note": "invalid username",
		})
	}
	password := c.FormValue("password")
	if password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"note": "invalid password",
		})
	}
	err := u.UserService.Register(username, password, c.Context())
	switch err {
	case entities.ErrUserAlreadyExists:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"note": "user already exists",
		})
	case nil:
		return c.SendStatus(fiber.StatusOK)
	default:
		return err
	}
}
