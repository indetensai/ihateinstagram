package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func follow_handler(c *fiber.Ctx) error {
	session_user_id := check_session(c)
	if session_user_id == nil {
		return nil
	}
	user_id := uuid.Must(uuid.Parse(c.Params("user_id")))
	follower_id := uuid.Must(uuid.Parse(c.Params("follower_id")))
	if follower_id != *session_user_id {
		return c.SendStatus(fiber.StatusForbidden)
	}
	_, err := con.Exec(c.Context(), "INSERT INTO following (follower_id,user_id) VALUES ($1,$2)", user_id, follower_id)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.SendStatus(fiber.StatusOK)
}
