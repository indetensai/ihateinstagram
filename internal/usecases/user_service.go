package usecases

import (
	"context"
	"crypto/rand"
	"math/big"
	"phota/internal/entities"

	"github.com/google/uuid"
)

type userService struct {
	userRepo entities.UserRepository
}

func NewUserService(u entities.UserRepository) entities.UserService {
	return &userService{u}
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
	user, err := u.userRepo.GetUserBySession(session_id, ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userService) DeleteSession(session_id string, ctx context.Context) error {
	err := u.userRepo.DeleteSession(session_id, ctx)
	return err
}

func (u *userService) Login(username string, password string, ctx context.Context) (*string, error) {
	session_id, err := GenerateRandomString(50)
	if err != nil {
		return nil, err
	}
	err = u.userRepo.ValidateUserCredentials(username, password, session_id, ctx)
	if err != nil {
		return nil, err
	}

	return &session_id, nil
}

func (u *userService) Register(username string, password string, ctx context.Context) error {
	return u.userRepo.CreateUser(username, password, ctx)
}
