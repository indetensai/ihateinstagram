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
	_, err := con.Exec(
		c.Context(),
		"INSERT INTO following (follower_id,user_id) VALUES ($1,$2)",
		follower_id,
		user_id,
	)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.SendStatus(fiber.StatusOK)
}
func unfollow_handler(c *fiber.Ctx) error {
	session_user_id := check_session(c)
	if session_user_id == nil {
		return nil
	}
	user_id := uuid.Must(uuid.Parse(c.Params("user_id")))
	follower_id := uuid.Must(uuid.Parse(c.Params("follower_id")))
	if follower_id != *session_user_id {
		return c.SendStatus(fiber.StatusForbidden)
	}
	_, err := con.Exec(
		c.Context(),
		"DELETE FROM following WHERE user_id=$1 AND follower_id=$2",
		user_id,
		follower_id,
	)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.SendStatus(fiber.StatusOK)
}
func followers_handler(c *fiber.Ctx) error {
	user_id := uuid.Must(uuid.Parse(c.Params("user_id")))
	rows, err := con.Query(
		c.Context(),
		"SELECT follower_id FROM following WHERE user_id=$1",
		user_id,
	)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	defer rows.Close()
	var followers []uuid.UUID
	for rows.Next() {
		var r uuid.UUID
		err = rows.Scan(&r)
		if err != nil {
			return err
		}
		followers = append(followers, r)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"followers": followers})
}
