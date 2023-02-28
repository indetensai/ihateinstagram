package entities

import "errors"

var (
	ErrEmptySession       = errors.New("empty session")
	ErrUserAlreadyExists  = errors.New("username already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrPostNotFound       = errors.New("Post not found")
	ErrNotAuthorized      = errors.New("You don't have permission to view this post")
)
