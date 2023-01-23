package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var con *pgx.Conn

func error_check(err error, status string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s - %v\n", status, err)
		os.Exit(1)
	}
}

func login_handler(c *fiber.Ctx) error {
	hash := sha256.New()
	hash.Write([]byte(c.FormValue("password")))
	pede := hash.Sum(nil)
	var b string
	var id uuid.UUID
	err := con.QueryRow(context.Background(), "SELECT password,user_id FROM app_user WHERE username=$1", c.FormValue("username")).Scan(&b, &id)
	if err != nil {
		fmt.Println(err)
		c.Write([]byte("username not found"))
		c.Status(fiber.StatusNotFound)
		return err
	}

	error_check(err, "scanning failed")

	if b == fmt.Sprintf("%x", pede) {
		_, err := con.Exec(context.Background(), "INSERT INTO session (user_id) VALUES ($1)", id)
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"note":    "sucessfull login!! congratz.",
		})
	} else {
		error_check(err, "response writing failed")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"note":    "password is wrong, ure brainless",
		})
	}
}

func register_handler(c *fiber.Ctx) error {
	hash := sha256.New()
	hash.Write([]byte(c.FormValue("password")))
	pede := hash.Sum(nil)
	var exists bool
	err := con.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM app_user WHERE username=$1)", c.FormValue("username")).Scan(&exists)
	error_check(err, "existence check failed")
	if exists || (c.FormValue("username") == "" || c.FormValue("password") == "") {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"note":    "username already exists or username or password blank",
		})

	} else {
		_, err := con.Exec(context.Background(), "INSERT INTO app_user (username,password) VALUES ($1,$2)", c.FormValue("username"), fmt.Sprintf("%x", pede))
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"note":    "",
		})
	}
}

func based() {
	var err error
	con, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	error_check(err, "connection to database failed")
	con.Ping(context.Background())
	con.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS app_user(
		user_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(64) NOT NULL
		);`)
	con.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS session(
		session_id UUID PRIMARY KEY NOT NULL UNIQUE DEFAULT gen_random_uuid() ,
		user_id uuid REFERENCES app_user(user_id) ON DELETE CASCADE,
		created_at TIMESTAMP NOT NULL DEFAULT now()	
		); `)
}
func main() {
	godotenv.Load(".env")
	based()
	app := fiber.New()
	app.Post("/user/register", register_handler)
	app.Post("/user/login", login_handler)
	app.Listen(":8080")
}
