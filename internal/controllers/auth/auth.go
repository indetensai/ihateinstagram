package auth

import (
	"phota/internal/entities"

	"github.com/gofiber/fiber/v2"
)

func New(repo entities.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session_id := c.Query("session_id")
		if session_id == "" {
			c.Locals("user_id", nil)
			return c.Next()
		}
		user_id, err := repo.GetUserBySession(session_id, c.Context())
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("invalid session")
		}
		c.Locals("user_id", user_id)
		return c.Next()
	}
}
