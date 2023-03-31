package app

import (
	"context"
	"log"
	"phota/internal/controllers"
	"phota/internal/controllers/auth"
	"phota/internal/repository"
	"phota/internal/usecases"

	"github.com/gofiber/fiber/v2"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

func Run(databaseURL string, port string) {
	con, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	if err = con.Ping(context.Background()); err != nil {
		log.Print("warn - database isn't respodning to ping: ", err)
	}

	app := fiber.New()

	user_repository := repository.NewUserRepository(con)
	user_service := usecases.NewUserService(user_repository)
	controllers.NewUserServiceHandler(app, user_service)

	app.Use(auth.New(user_repository))

	following_repository := repository.NewFollowingRepository(con)
	following_service := usecases.NewFollowingService(following_repository, user_service)
	controllers.NewFollowingServiceHandler(app, following_service)

	post_repository := repository.NewPostRepository(con)
	post_service := usecases.NewPostService(post_repository, user_service, following_service)
	controllers.NewPostServiceHandler(app, post_service, following_service)

	image_repository := repository.NewImageRepository(con)
	image_service := usecases.NewImageService(image_repository, post_service, following_service)
	controllers.NewImageServiceHandler(app, image_service)

	app.Listen(":" + port)
}
