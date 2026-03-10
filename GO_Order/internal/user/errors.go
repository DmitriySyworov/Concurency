package user

import "errors"

var (
	ErrReg = errors.New("this user already exists, please log in")

	ErrSecurity = errors.New("failed to ensure the security of data transmission")

	ErrSess = errors.New("the authorization session has expired or is not valid")

	ErrSendEmail = errors.New("we were unable to send an email to the specified email address")

	ErrMissing    = errors.New("the required data has not been transmitted")
	ErrMethod     = errors.New("there is no such authorization method")
	ErrCreateUser = errors.New("failed to records user")
	ErrWrongData  = errors.New("wrong email/phone or password")
	ErrAuth       = errors.New("there is no such authorization method")

	ErrUserNotFound = errors.New("such user not found")

	ErrIncorrectCode = errors.New("code is incorrect")
)
