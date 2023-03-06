package controllers

import (
	"phota/internal/entities"

	"github.com/gofiber/fiber/v2"
)

type userServiceHandler struct {
	UserService entities.UserService
}

func NewUserServiceHandler(app *fiber.App, u entities.UserService) {
	handler := &userServiceHandler{u}
	app.Post("/user/register", handler.RegisterHandler)
	app.Post("/user/login", handler.LoginHandler)
	app.Delete("/session/:id<guid>", handler.DeleteSessionHandler)
}

func (u *userServiceHandler) RegisterHandler(c *fiber.Ctx) error {
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
	if err != nil {
		return error_handling(c, err)
	}
	return nil
}
func (u *userServiceHandler) LoginHandler(c *fiber.Ctx) error {
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
	session_id, err := u.UserService.Login(username, password, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"session_id": session_id})
}
func (u *userServiceHandler) DeleteSessionHandler(c *fiber.Ctx) error {
	session_id := c.FormValue("session_id")
	err := u.UserService.DeleteSession(session_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}
