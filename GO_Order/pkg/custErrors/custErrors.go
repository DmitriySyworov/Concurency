package custerrors

import "errors"

var (
	ErrResponse = errors.New("failed to generate response")
	ErrRequest = errors.New("incorrect request structure")
	ErrInvalidData = errors.New("invalid data transmitted")
	ErrInvalidToken = errors.New("you passed the wrong token")
	ErrToken = errors.New("Failed to securely write transferred data")
)