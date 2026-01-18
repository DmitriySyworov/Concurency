package payload

import "math/rand/v2"

type Verification struct {
	Email string `json:"email" validate:"required,email"`
	Hash  string `json:"hash"`
}
type ResponseSend struct {
	Send  bool   `json:"send"`
	Error string `json:"error"`
}
type ResponseVerify struct{
	Success  bool   `json:"success"`
	Error string `json:"error"`
}
func NewVerification(email string) *Verification {
	return &Verification{
		Email: email,
		Hash:  generateLetter(),
	}
}
func generateLetter() string {
	i := 0
	resLetter := make([]byte, 0, 9)
	for 9 > i {
		randomNum := rand.IntN(123)
		if (randomNum > 47 && randomNum < 58) || (randomNum > 64 && randomNum < 91) || (randomNum > 96) {
			resLetter = append(resLetter, byte(randomNum))
			i++
		}
	}
	return string(resLetter)
}
func NewResponseSend(err error) *ResponseSend {
	if err != nil {
		return &ResponseSend{
			Error: err.Error(),
			Send:  false,
		}
	}
	return &ResponseSend{
		Error: "",
		Send:  true,
	}
}
func NewResponseVerify(err error)*ResponseVerify{
	if err != nil {
		return &ResponseVerify{
			Error: err.Error(),
			Success:  false,
		}
	}
	return &ResponseVerify{
		Error: "",
		Success:  true,
	}
}