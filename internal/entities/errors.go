package entities

import "errors"

var (
	ErrEmptySession       = errors.New("empty session")
	ErrUserAlreadyExists  = errors.New("username already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
)
