package app

import (
	"context"
	"log"
	"os"
	"phota/internal/controllers"
	"phota/internal/controllers/auth"
	"phota/internal/repository"
	"phota/internal/usecases"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/skip"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func Run() {
	godotenv.Load(".env")
	m, err := migrate.New("file://migrations", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("failed to read migration ", err)
	}
	if err := m.Up(); err != nil {
		log.Print("warn - migration failed: ", err)
	}
	con, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
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

	app.Use(skip.New(auth.New(user_repository), func(c *fiber.Ctx) bool {
		return c.Path() == "/user/register" || c.Path() == "/user/login" || c.Path() == "/user/:user_id<guid>/followers"
	}))

	following_repository := repository.NewFollowingRepository(con)
	following_service := usecases.NewFollowingService(following_repository, user_service)
	controllers.NewFollowingServiceHandler(app, following_service)

	post_repository := repository.NewPostRepository(con)
	post_service := usecases.NewPostService(post_repository, user_service)
	controllers.NewPostServiceHandler(app, post_service, following_service)

	image_repository := repository.NewImageRepository(con)
	image_service := usecases.NewImageService(image_repository)
	controllers.NewImageServiceHandler(app, image_service, post_service, following_service)

	app.Listen(":8080")
}
