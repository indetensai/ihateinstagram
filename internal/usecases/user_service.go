package usecases

import (
	"context"
	"crypto/rand"
	"math/big"
	"phota/internal/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	db *pgx.Conn
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-."
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "generating failed", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func (u *userService) CheckSession(session_id string, ctx context.Context) (*uuid.UUID, error) {
	if session_id == "" {
		return nil, entities.ErrEmptySession
	}
	var user_id *uuid.UUID
	err := u.db.QueryRow(
		ctx,
		"SELECT user_id FROM session WHERE session_id=$1",
		session_id,
	).Scan(user_id)
	if err != nil {
		return nil, err
	}
	return user_id, nil
}

func (u *userService) DeleteSession(session_id string, ctx context.Context) error {
	_, err := u.db.Exec(
		ctx,
		"DELETE FROM session WHERE session_id=$1",
		session_id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (u *userService) Login(username string, password string, ctx context.Context) error {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var result entities.AppUser
	err = tx.QueryRow(
		ctx,
		"SELECT * FROM app_user WHERE username=$1",
		username,
	).Scan(&result)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password))
	if err != nil {
		return entities.ErrInvalidCredentials
	}
	session_id, err := GenerateRandomString(64)
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		ctx,
		"INSERT INTO session (session_id,user_id) VALUES ($1,$2)",
		session_id, result.UserID,
	)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (u *userService) Register(username string, password string, ctx context.Context) error {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	password_raw := []byte(password)
	password_hashed, err := bcrypt.GenerateFromPassword(password_raw, 10)
	if err != nil {
		return err
	}
	var exists bool
	err = tx.QueryRow(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM app_user WHERE username=$1)",
		username,
	).Scan(&exists)
	if err != nil {
		return err
	}
	if exists || (username == "" || password == "") {
		return entities.ErrUserAlreadyExists

	} else {
		_, err := tx.Exec(
			ctx,
			"INSERT INTO app_user (username,password) VALUES ($1,$2)",
			username,
			string(password_hashed),
		)
		if err != nil {
			return err
		}
		tx.Commit(ctx)
		return nil
	}
}
