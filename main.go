package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
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
	app.Delete("/session/:id<guid>", delete_session)
	app.Listen(":8080")
}
