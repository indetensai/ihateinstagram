package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var con *pgx.Conn

type register_result struct {
	Success bool   `json:"success"`
	Note    string `json:"note,omitempty"`
}

func error_check(err error, status string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s - %v\n", status, err)
		os.Exit(1)
	}
}

func register_handler(c *fiber.Ctx) error {
	id := uuid.New()
	hash := sha256.New()
	hash.Write([]byte(c.FormValue("password")))
	pede := hash.Sum(nil)
	var exists bool
	letter := new(register_result)
	err := con.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM app_user WHERE username=$1)", c.FormValue("username")).Scan(&exists)
	error_check(err, "existence check failed")
	if exists || (c.FormValue("username") == "" || c.FormValue("password") == "") {
		letter.Success = false
		letter.Note = "username already exists or username or password blank"

	} else if !exists {
		letter.Success = true
		letter.Note = ""
		_, err := con.Exec(context.Background(), "INSERT INTO app_user (user_id,username,password) VALUES ($1,$2,$3)", id, c.FormValue("username"), fmt.Sprintf("%x", pede))
		error_check(err, "inserting new user credentials failed")
	}
	j, err := json.Marshal(letter)
	error_check(err, "response writing failed")
	c.Write(j)
	return err
}

func baza() {
	var err error
	con, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	error_check(err, "connection to database failed")
	con.Ping(context.Background())
	con.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS app_user(
		user_id uuid PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(64) NOT NULL
		);`)
}
func main() {
	godotenv.Load(".env")
	baza()
	app := fiber.New()
	app.Post("/user/register", register_handler)
	app.Listen(":8080")
}
