package custerrors

import "errors"

var (
	ErrInvalidToken = errors.New("invalid or expired token")

	ErrResponse    = errors.New("failed to generate response")
	ErrRequest     = errors.New("incorrect requestJs structure")
	ErrInvalidData = errors.New("invalid data transmitted")

	ErrUserDontExist = errors.New("such auth does not exist")

	ErrCreateToken = errors.New("failed to create token")

	ErrUserNotFound = errors.New("such user not found")
)
