package jwt

import (
	custerrors "order/app/pkg/custErrors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	Secret []byte
}
type DataJWt struct {
	Email  string
	Phone  string
	IdUser string
}

func NewJWT(secret []byte) *JWT {
	return &JWT{
		Secret: secret,
	}
}
func (j *JWT) CreateJWT(data *DataJWt) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":     data.Email,
		"phone":     data.Phone,
		"id_user":   data.IdUser,
		"temporary": false,
	})
	resJWT, errToken := token.SignedString(j.Secret)
	if errToken != nil {
		return "", errToken
	}
	return resJWT, nil
}

func (j *JWT) TemporaryJWT(data *DataJWt) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":     data.Email,
		"phone":     data.Phone,
		"id_user":   data.IdUser,
		"temporary": true,
		"ExpiresAt": time.Now().Add(5 * time.Minute).Unix(),
	})
	resJWT, errToken := token.SignedString(j.Secret)
	if errToken != nil {
		return "", errToken
	}
	return resJWT, nil
}
func (j *JWT) DecodeJWT(tokenUs string) (*DataJWt, error) {
	token, errToken := jwt.Parse(tokenUs, func(t *jwt.Token) (any, error) {
		return j.Secret, nil
	})
	if errToken != nil || !token.Valid {
		return nil, custerrors.ErrInvalidToken
	}
	email, okEm := token.Claims.(jwt.MapClaims)["email"].(string)
	phone, okPh := token.Claims.(jwt.MapClaims)["phone"].(string)
	id, okId := token.Claims.(jwt.MapClaims)["id_user"].(string)
	if !okId && okPh && okEm {
		return &DataJWt{Email: email, Phone: phone}, nil
	}
	if okId && !okPh && !okEm {
		return &DataJWt{IdUser: id}, nil
	}
	return &DataJWt{Email: email, Phone: phone, IdUser: id}, nil
}
