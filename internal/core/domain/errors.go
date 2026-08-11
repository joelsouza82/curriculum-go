package domain

import "errors"

var (
	ErrPersonalNotFound = errors.New("personal not found")
	ErrLoginNotFound    = errors.New("login not found")
)
