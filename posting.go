package main

import (
	b64 "encoding/base64"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type post struct {
	owner_id    uuid.UUID
	description string
	visibility  string
	created_at  time.Time
}

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
	_, err = tx.Exec(
		c.Context(),
		"INSERT INTO post (user_id,description) VALUES ($1,$2)",
		user_id,
		c.FormValue("description"),
	)
	if err != nil {
		return err
	}
	var session_id uuid.UUID
	err = tx.QueryRow(
		c.Context(),
		"SELECT post_id FROM post ORDER BY created_at DESC",
	).Scan(&session_id)
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
	err = tx.QueryRow(
		c.Context(),
		"SELECT user_id FROM post WHERE post_id=$1",
		post_id,
	).Scan(&check)
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
	_, err = tx.Exec(
		c.Context(),
		"INSERT INTO image (user_id,post_id,content) VALUES ($1,$2,$3)",
		user_id,
		post_id,
		image,
	)
	if err != nil {
		return err
	}
	tx.Commit(c.Context())
	return c.SendStatus(fiber.StatusOK)
}
func getting_post_handler(c *fiber.Ctx) error {
	post := new(post)
	user_id := check_session(c)
	post_id := c.Params("post_id")
	tx, err := con.Begin(c.Context())
	if err != nil {
		return err
	}
	defer tx.Rollback(c.Context())
	err = tx.QueryRow(
		c.Context(),
		"SELECT user_id,description,vision,created_at FROM post WHERE post_id=$1",
		post_id,
	).Scan(
		&post.owner_id,
		&post.description,
		&post.visibility,
		&post.created_at,
	)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if post.visibility == "private" {
		if user_id == nil || *user_id != post.owner_id {
			return c.SendStatus(fiber.StatusNotFound)
		}
	}
	if post.visibility == "followers" {
		if user_id != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		_, err = tx.Exec(
			c.Context(),
			"SELECT 1 FROM following WHERE follower_id=$1 AND user_id=$2",
			user_id,
			post.owner_id,
		)
		if err != nil && *user_id != post.owner_id {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}

	tx.Commit(c.Context())
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"post_owner_id": post.owner_id,
			"visibility":    post.visibility,
			"created_at":    post.created_at,
			"description":   post.description,
		},
	)
}

func post_changing_handler(c *fiber.Ctx) error {
	user_id := check_session(c)
	if user_id == nil {
		return nil
	}
	received_post := new(post)
	post_id := c.Params("post_id")
	tx, err := con.Begin(c.Context())
	if err != nil {
		return err
	}
	defer tx.Rollback(c.Context())
	err = tx.QueryRow(
		c.Context(),
		"SELECT user_id,description,vision,created_at FROM post WHERE post_id=$1",
		post_id,
	).Scan(
		&received_post.owner_id,
		&received_post.description,
		&received_post.visibility,
		&received_post.created_at,
	)
	if *user_id != received_post.owner_id {
		c.SendStatus(fiber.StatusUnauthorized)
	}
	changes := struct {
		Visibility  string `json:"visibility"`
		Description string `json:"description"`
	}{}
	check := *received_post
	if err := c.BodyParser(&changes); err != nil {
		return err
	}
	if changes.Description != check.description && changes.Description != "" {
		check.description = changes.Description
	}
	if changes.Visibility != check.visibility && changes.Visibility != "" {
		check.visibility = changes.Visibility
	}
	if check != *received_post {
		_, err = tx.Exec(
			c.Context(),
			"UPDATE post SET vision=$1,description=$2 WHERE post_id=$3",
			check.visibility,
			check.description,
			post_id,
		)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
	}
	tx.Commit(c.Context())
	return c.SendStatus(fiber.StatusOK)
}
func like_handler(c *fiber.Ctx) error {
	tx, err := con.Begin(c.Context())
	if err != nil {
		return err
	}
	defer tx.Rollback(c.Context())
	user_id := check_session(c)
	if user_id == nil {
		return nil
	}
	post_id := c.Params("post_id")
	_, err = tx.Exec(
		c.Context(),
		"SELECT 1 FROM post WHERE post_id=$1",
		post_id,
	)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	_, err = tx.Exec(
		c.Context(),
		"INSERT INTO liking (user_id,post_id) VALUES ($1,$2)",
		user_id,
		post_id,
	)
	if err != nil {
		return c.SendStatus(fiber.StatusConflict)
	}
	tx.Commit(c.Context())
	return c.SendStatus(fiber.StatusOK)
}
