package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	database_connection()
	app := fiber.New()
	app.Post("/user/register", register_handler)
	app.Post("/user/login", login_handler)
	app.Delete("/session/:id<guid>", delete_session)
	app.Post("/post", post_handler)
	app.Listen(":8080")
}
