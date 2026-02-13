package user

import "errors"
var (
	ErrSecurity = errors.New("failed to ensure the security of data transmission")

	ErrSess = errors.New("the authorization session has expired or is not valid")

	ErrMissing = errors.New("the required data has not been transmitted")
	ErrMethod = errors.New("there is no such authorization method")
	ErrCreateUser = errors.New("failed to records user")
	ErrWrongData = errors.New("wrong email/phone or password")
	ErrAuth = errors.New("there is no such authorization method")
	ErrCode = errors.New("the code has expired")
	ErrIncorrectCode = errors.New("code is incorrect")
)