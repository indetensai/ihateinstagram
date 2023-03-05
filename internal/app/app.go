package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"phota/internal/controllers"
	"phota/internal/usecases"
	"phota/internal/usecases/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func Run() {
	godotenv.Load(".env")
	m, err := migrate.New("file://../sql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
	con, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Print(err)
		log.Fatal("baza caput")
	}
	err = con.Ping(context.Background())
	if err != nil {
		fmt.Print(err)
		log.Fatal("baza caput")
	}

	app := fiber.New()

	user_repository := repository.NewUserRepository(con)
	user_service := usecases.NewUserService(user_repository)
	controllers.NewUserServiceHandler(app, user_service)

	following_repository := repository.NewFollowingRepository(con)
	following_service := usecases.NewFollowingService(following_repository, user_service)
	controllers.NewFollowingServiceHandler(app, following_service)

	post_repository := repository.NewPostRepository(con)
	post_service := usecases.NewPostService(post_repository, user_service)
	controllers.NewPostServiceHandler(app, post_service)

	app.Listen(":8080")
}
