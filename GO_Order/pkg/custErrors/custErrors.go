package custerrors

import "errors"

var (
	ErrInvalidToken = errors.New("invalid or expired token")

	ErrResponse    = errors.New("failed to generate response")
	ErrRequest     = errors.New("incorrect request structure")
	ErrInvalidData = errors.New("invalid data transmitted")

	ErrUserDontExist = errors.New("such user does not exist")

	ErrCreateToken = errors.New("failed to create token")
)
