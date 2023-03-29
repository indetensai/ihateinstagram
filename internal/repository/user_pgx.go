package repository

import (
	"context"
	"errors"
	"phota/internal/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) entities.UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) GetUserBySession(session_id string, ctx context.Context) (*uuid.UUID, error) {
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

func (u *UserRepository) ValidateUserCredentials(
	username string,
	password string,
	session_id string,
	ctx context.Context,
) error {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var user entities.UserCredentials
	err = tx.QueryRow(
		ctx,
		"SELECT * FROM app_user WHERE username=$1",
		username,
	).Scan(&user.UserID, &user.Username, &user.Password)
	if err != nil {
		return entities.ErrNotFound
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return entities.ErrInvalidCredentials
	}
	_, err = tx.Exec(
		ctx,
		"INSERT INTO session (session_id,user_id) VALUES ($1,$2)",
		session_id, user.UserID,
	)
	if err != nil {
		return entities.ErrDuplicate
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) CreateUser(username string, password string, ctx context.Context) error {
	password_raw := []byte(password)
	password_hashed, err := bcrypt.GenerateFromPassword(password_raw, 10)
	if err != nil {
		return err
	}
	_, err = u.db.Exec(
		ctx,
		"INSERT INTO app_user (username,password) VALUES ($1,$2)",
		username,
		string(password_hashed),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			return entities.ErrDuplicate
		}
		return err
	}
	return nil
}

func (u *UserRepository) GetUserByID(user_id uuid.UUID, ctx context.Context) (*entities.AppUser, error) {
	var user entities.AppUser
	err := u.db.QueryRow(
		ctx,
		"SELECT (user_id,username) FROM app_user WHERE user_id=$1",
		user_id,
	).Scan(&user)
	if err != nil {
		return nil, entities.ErrNotFound
	}
	return &user, nil
}
