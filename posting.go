package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func post_handler(c *fiber.Ctx) error {
	user_id := check_session(c)
	if user_id == nil {
		return nil
	}
	tx, err := con.Begin(c.Context())
	if err != nil {
		return err
	}
	defer tx.Rollback(c.Context())
	_, err = tx.Exec(c.Context(), "INSERT INTO post (user_id,description) VALUES ($1,$2)", user_id, c.FormValue("description"))
	if err != nil {
		return err
	}
	var session_id uuid.UUID
	err = tx.QueryRow(c.Context(), "SELECT post_id FROM post ORDER BY created_at DESC").Scan(&session_id)
	if err != nil {
		return err
	}
	err = tx.Commit(c.Context())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"post_id": session_id,
	})
}
