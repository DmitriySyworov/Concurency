package custerrors

import "errors"

var (
	ErrInvalidToken = errors.New("invalid or expired token")

	ErrResponse    = errors.New("failed to generate response")
	ErrRequest     = errors.New("incorrect request structure")
	ErrInvalidData = errors.New("invalid data transmitted")

	ErrCreateToken = errors.New("failed to create token")
)
