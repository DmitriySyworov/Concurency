package auth

type RequestUserRegister struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type RequestUserLogin struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required_without=Phone"`
	Phone    string `json:"phone" validate:"required_without=Email"`
	Password string `json:"password" validate:"required"`
}
type RequestRestore struct {
	Email string `json:"email" validate:"required_without=Phone"`
	Phone string `json:"phone" validate:"required_without=Email"`
}
type ResponseAuth struct {
	Message   string `json:"message"`
	SessionId string `json:"session_id"`
	Jwt       string `json:"jwt"`
	Error     string `json:"error"`
}
type RequestConfirm struct {
	TempPassword int `json:"tempPassword"`
}
type ResponseConfirm struct {
	Jwt   string `json:"jwt"`
	Error string `json:"error"`
}
