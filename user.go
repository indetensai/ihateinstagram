package main

import (
	"crypto/sha256"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func delete_session(c *fiber.Ctx) error {
	_, err := con.Exec(c.Context(), "DELETE FROM session WHERE session_id=$1", c.Params("id"))
	if err != nil {
		return err
	} else {
		c.Status(fiber.StatusOK)
	}
	return nil
}

func login_handler(c *fiber.Ctx) error {
	tx, err := con.Begin(c.Context())
	if err != nil {
		return err
	}
	defer tx.Rollback(c.Context())
	hash := sha256.New()
	hash.Write([]byte(c.FormValue("password")))
	pede := hash.Sum(nil)
	var b string
	var id uuid.UUID
	err = tx.QueryRow(c.Context(), "SELECT password,user_id FROM app_user WHERE username=$1", c.FormValue("username")).Scan(&b, &id)
	if err != nil {
		fmt.Println(err)
		c.Write([]byte("username not found"))
		c.Status(fiber.StatusNotFound)
		return err
	}

	error_check(err, "scanning failed")

	if b == fmt.Sprintf("%x", pede) {
		_, err := tx.Exec(c.Context(), "INSERT INTO session (user_id) VALUES ($1)", id)
		if err != nil {
			return err
		}
		var session_id string
		err = tx.QueryRow(c.Context(), "SELECT session_id FROM session ORDER BY created_at DESC ").Scan(&session_id)
		if err != nil {
			return err
		}
		err = tx.Commit(c.Context())
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"session": session_id,
			"note":    "sucessfull login!! congratz.",
		})
	} else {
		error_check(err, "response writing failed")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"note": "password is wrong, ure brainless",
		})
	}
}

func register_handler(c *fiber.Ctx) error {
	tx, err := con.Begin(c.Context())
	if err != nil {
		return err
	}
	defer tx.Rollback(c.Context())
	hash := sha256.New()
	hash.Write([]byte(c.FormValue("password")))
	pede := hash.Sum(nil)
	var exists bool
	err = tx.QueryRow(c.Context(), "SELECT EXISTS(SELECT 1 FROM app_user WHERE username=$1)", c.FormValue("username")).Scan(&exists)
	error_check(err, "existence check failed")
	if exists || (c.FormValue("username") == "" || c.FormValue("password") == "") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"note": "username already exists or username or password blank",
		})

	} else {
		_, err := tx.Exec(c.Context(), "INSERT INTO app_user (username,password) VALUES ($1,$2)", c.FormValue("username"), fmt.Sprintf("%x", pede))
		if err != nil {
			return err
		}
		tx.Commit(c.Context())
		return c.Status(fiber.StatusOK).JSON(fiber.Map{})
	}
}
