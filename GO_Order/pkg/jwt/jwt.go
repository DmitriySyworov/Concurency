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
	IdUser float64
}

func NewJWT(secret []byte) *JWT {
	return &JWT{
		Secret: secret,
	}
}
func (j *JWT) CreateJWT(data *DataJWt) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id_user": data.IdUser,
	})
	resJWT, errToken := token.SignedString(j.Secret)
	if errToken != nil {
		return "", custerrors.ErrCreateToken
	}
	return resJWT, nil
}

func (j *JWT) TemporaryJWT(data *DataJWt) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":     data.Email,
		"ExpiresAt": time.Now().Add(5 * time.Minute).Unix(),
	})
	resJWT, errToken := token.SignedString(j.Secret)
	if errToken != nil {
		return "", custerrors.ErrCreateToken
	}
	return resJWT, nil
}
func (j *JWT) ParseJWT(tokenUs string) (*DataJWt, error) {
	token, errToken := jwt.Parse(tokenUs, func(t *jwt.Token) (any, error) {
		return j.Secret, nil
	})
	if errToken != nil || !token.Valid {
		return nil, custerrors.ErrInvalidToken
	}
	if id, okId := token.Claims.(jwt.MapClaims)["id_user"].(float64); okId {
		return &DataJWt{IdUser: id}, nil
	}
	return nil, custerrors.ErrInvalidToken
}
func (j *JWT) ParseTemporaryJWT(tokenUs string) (*DataJWt, error) {
	token, errToken := jwt.Parse(tokenUs, func(t *jwt.Token) (any, error) {
		return j.Secret, nil
	})
	if errToken != nil || !token.Valid {
		return nil, custerrors.ErrInvalidToken
	}
	email, okEm := token.Claims.(jwt.MapClaims)["email"].(string)
	if !okEm {
		return nil, custerrors.ErrInvalidToken
	}
	return &DataJWt{Email: email}, nil
}
