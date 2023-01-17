package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber"
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

/*func register_handler(fiber.Ctx) {
	id := uuid.New()
	hash := sha256.New()
	hash.Write([]byte(r.FormValue("password")))
	pede := hash.Sum(nil)
	var exists bool
	letter := new(register_result)
	err := con.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM app_user WHERE username=$1)", r.FormValue("username")).Scan(&exists)
}*/

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
	//app.Post("/register")

}
