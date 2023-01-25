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
		"CREATE TYPE visibility AS ENUM ('public', 'followers', 'private');",
	)
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
		description TEXT,
		vision visibility DEFAULT 'private',
		created_at TIMESTAMP NOT NULL DEFAULT now()
		);`)
	con.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS image(
		image_id UUID PRIMARY KEY UNIQUE DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES app_user(user_id) ON DELETE CASCADE,
		post_id UUID REFERENCES post(post_id) ON DELETE SET NULL,
		content BYTEA 
		);`)
	con.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS following(
		follower_id UUID PRIMARY KEY NOT NULL REFERENCES app_user(user_id) ON DELETE CASCADE,
		user_id UUID REFERENCES app_user(user_id) ON DELETE CASCADE,
		following_since TIMESTAMP NOT NULL DEFAULT now(),
		UNIQUE (user_id,follower_id),
		CHECK (user_id!=follower_id)
		);`)
	con.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS liking(
		user_id UUID NOT NULL REFERENCES app_user(user_id) ON DELETE CASCADE,
		post_id UUID NOT NULL REFERENCES post(post_id) ON DELETE CASCADE,
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		PRIMARY KEY (user_id,post_id),
		UNIQUE (user_id,post_id)
		);`)
}
