package auth

import "errors"

var (
	ErrReg = errors.New("this auth already exists, please log in")

	ErrSecurity = errors.New("failed to ensure the security of data transmission")

	ErrSess = errors.New("the authorization session has expired or is not valid")

	ErrSendEmail = errors.New("we were unable to send an email to the specified email address")

	ErrMissing    = errors.New("the required data has not been transmitted")
	ErrCreateUser = errors.New("failed to records auth")
	ErrWrongData  = errors.New("wrong email/phone or password")
	ErrMethodAuth = errors.New("there is no such authorization method or the data does not match the selected method")

	ErrIncorrectCode = errors.New("code is incorrect")

	ErrRestoreUser = errors.New("failed to restore user")
)
