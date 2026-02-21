package user

import "time"

type RequestUserRegist struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type RequestUserLogin struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password" validate:"required"`
}
type ResponseAuth struct {
	SessionId string `json:"sessionid"`
	Message   string `json:"message"`
	Error     string `json:"error"`
}
type RequestConfirm struct {
	TempPassword string `json:"tempPassword"`
}
type ResponseConfirm struct {
	Jwt   string `json:"jwt"`
	Error string `json:"error"`
}
type TempJWTUser struct {
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Temporary bool`json:"temporary"`
	ExpiresAt time.Time `json:"ExpiresAr"`
}
