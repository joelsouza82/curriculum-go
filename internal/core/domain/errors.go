package domain

import "errors"

var (
	ErrPersonalNotFound   = errors.New("personal not found")
	ErrLoginNotFound      = errors.New("login not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
)
