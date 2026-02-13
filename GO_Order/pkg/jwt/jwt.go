package jwt

import (
	custerrors "order/app/pkg/custErrors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	Secret []byte
}

func NewJWT(secret []byte) *JWT {
	return &JWT{
		Secret: secret,
	}
}
func (j *JWT) CreateJWT(email, phone string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"phone": phone,
	})
	resJWT, errToken := token.SignedString(j.Secret)
	if errToken != nil {
		return "", errToken
	}
	return resJWT, nil
}

func (j *JWT) FastJWT(v any) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":     v,
		"issuedAt": jwt.NewNumericDate(time.Now()),
	})
	resJWT, errToken := token.SignedString(j.Secret)
	if errToken != nil {
		return "", errToken
	}
	return resJWT, nil
}
func (j *JWT) DecodeJWT(tokenUs string) (jwt.MapClaims, error) {
	token, errToken := jwt.Parse(tokenUs, func(t *jwt.Token) (any, error) {
		return j.Secret, nil
	})
	if errToken != nil {
		return nil, custerrors.ErrInvalidToken
	}
	return token.Claims.(jwt.MapClaims), nil
}
