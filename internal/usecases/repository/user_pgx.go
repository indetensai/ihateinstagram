package repository

import (
	"context"
	"crypto/rand"
	"math/big"
	"phota/internal/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) entities.UserRepository {
	return &UserRepository{db: db}
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

func (u *UserRepository) CheckSession(session_id string, ctx context.Context) (*uuid.UUID, error) {
	var user_id uuid.UUID
	err := u.db.QueryRow(
		ctx,
		"SELECT user_id FROM session WHERE session_id=$1",
		session_id,
	).Scan(&user_id)
	if err != nil {
		return nil, entities.ErrNotFound
	}
	return &user_id, nil
}

func (u *UserRepository) DeleteSession(session_id string, ctx context.Context) error {
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

func (u *UserRepository) Login(username string, password string, ctx context.Context) (*string, error) {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	var result entities.AppUser
	err = tx.QueryRow(
		ctx,
		"SELECT * FROM app_user WHERE username=$1",
		username,
	).Scan(&result.UserID, &result.Name, &result.Password)
	if err != nil {
		return nil, entities.ErrNotFound
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password))
	if err != nil {
		return nil, entities.ErrInvalidCredentials
	}
	session_id, err := GenerateRandomString(64)
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(
		ctx,
		"INSERT INTO session (session_id,user_id) VALUES ($1,$2)",
		session_id, result.UserID,
	)
	if err != nil {
		return nil, entities.ErrDuplicate
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return &session_id, nil
}

func (u *UserRepository) Register(username string, password string, ctx context.Context) error {
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
		return entities.ErrNotFound
	}
	if exists || (username == "" || password == "") {
		return entities.ErrDuplicate

	} else {
		_, err := tx.Exec(
			ctx,
			"INSERT INTO app_user (username,password) VALUES ($1,$2)",
			username,
			string(password_hashed),
		)
		if err != nil {
			return entities.ErrDuplicate
		}
		tx.Commit(ctx)
		return nil
	}
}
