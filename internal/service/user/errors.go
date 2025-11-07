package user

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailInUse   = errors.New("email already in use")
	ErrInvalidInput = errors.New("invalid input: name and email are required")
)
