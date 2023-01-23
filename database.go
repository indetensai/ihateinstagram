package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var con *pgx.Conn

func error_check(err error, status string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s - %v\n", status, err)
		os.Exit(1)
	}
}

func database_connection() {
	var err error
	con, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	error_check(err, "connection to database failed")
	con.Ping(context.Background())
	con.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS app_user(
		user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(64) NOT NULL
		);`)
	con.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS session(
		session_id UUID PRIMARY KEY UNIQUE DEFAULT gen_random_uuid(),
		user_id UUID REFERENCES app_user(user_id) ON DELETE CASCADE,
		created_at TIMESTAMP NOT NULL DEFAULT now()	
		);`)
	con.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS post(
		post_id UUID PRIMARY KEY UNIQUE DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES app_user(user_id) ON DELETE CASCADE,
		description TEXT
		);`)
	con.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS image(
		image_id UUID PRIMARY KEY UNIQUE DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES app_user(user_id) ON DELETE CASCADE,
		post_id UUID REFERENCES post(post_id) ON DELETE SET NULL,
		content BYTEA 
		);`)
}
