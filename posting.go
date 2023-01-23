package main

import (
	b64 "encoding/base64"

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
func image_handler(c *fiber.Ctx) error {
	user_id := check_session(c)
	if user_id == nil {
		return nil
	}
	tx, err := con.Begin(c.Context())
	if err != nil {
		return err
	}
	defer tx.Rollback(c.Context())
	post_id := c.Params("post_id")
	var check uuid.UUID
	err = tx.QueryRow(c.Context(), "SELECT user_id FROM post WHERE post_id=$1", post_id).Scan(&check)
	if err != nil {
		return err
	}
	if check != *user_id {
		return c.SendStatus(fiber.StatusForbidden)
	}
	image_raw := c.Query("image")
	if image_raw == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	image, err := b64.StdEncoding.DecodeString(image_raw)
	if err != nil {
		return err
	}
	_, err = tx.Exec(c.Context(), "INSERT INTO image (user_id,post_id,content) VALUES ($1,$2,$3)", user_id, post_id, image)
	if err != nil {
		return err
	}
	tx.Commit(c.Context())
	return c.SendStatus(fiber.StatusOK)
}
