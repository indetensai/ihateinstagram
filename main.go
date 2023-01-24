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
	app.Put("/post/:post_id<guid>/image", image_handler)
	app.Put("/user/:user_id<guid>/followers/:follower_id<guid>", follow_handler)
	app.Delete("/user/:user_id<guid>/followers/:follower_id<guid>", unfollow_handler)
	app.Get("/user/:user_id<guid>/followers", followers_handler)
	app.Get("/post/:post_id<guid>", getting_post_handler)
	app.Listen(":8080")
}
