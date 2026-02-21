package custerrors

import "errors"

var (
	ErrInvalidToken = errors.New("Invalid or expired token")

	ErrResponse     = errors.New("failed to generate response")
	ErrRequest      = errors.New("incorrect request structure")
	ErrInvalidData  = errors.New("invalid data transmitted")
	ErrWrongJwt = errors.New("you passed the wrong token")
	ErrToken        = errors.New("Failed to securely write transferred data")
)
