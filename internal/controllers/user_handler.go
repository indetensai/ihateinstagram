package controllers

import (
	"phota/internal/entities"

	"github.com/gofiber/fiber/v2"
)

type UserServiceHandler struct {
	UserService entities.UserService
}

func NewUserServiceHandler(app *fiber.App, u entities.UserService) {
	handler := &UserServiceHandler{u}
	app.Post("/user/register", handler.RegisterHandler)
	app.Post("/user/login", handler.LoginHandler)
	app.Delete("/session/:id<guid>", handler.DeleteSessionHandler)
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
func (u *UserServiceHandler) LoginHandler(c *fiber.Ctx) error {
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
	if err == entities.ErrInvalidCredentials {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"session_id": session_id})
}
func (u *UserServiceHandler) DeleteSessionHandler(c *fiber.Ctx) error {
	session_id := c.FormValue("session_id")
	err := u.UserService.DeleteSession(session_id, c.Context())
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}
