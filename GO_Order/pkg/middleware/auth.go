package middleware

import (
	"context"
	"net/http"
	"order/app/configs"
	custerrors "order/app/pkg/custErrors"
	"order/app/pkg/jwt"
	"strings"
)

type KeyCtx string

const (
	KeyIDUser   KeyCtx = "KeyIDUser"
	KeyTempUser KeyCtx = "KeyTempUser"
)

func validateToken(header string) (string, error) {
	if !strings.HasPrefix(header, "Bearer ") {
		return "", custerrors.ErrInvalidToken
	}
	if strings.Count(header, ".") != 2 {
		return "", custerrors.ErrInvalidToken
	}
	resToken := strings.TrimPrefix(header, "Bearer ")
	if resToken == "" {
		return "", custerrors.ErrInvalidToken
	}
	return resToken, nil
}
func IsAuthID(next http.Handler, conf *configs.Config) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		token, errValid := validateToken(authHeader)
		if errValid != nil {
			http.Error(writer, custerrors.ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}
		parseJwt, errParse := jwt.NewJWT(conf.Secret).ParseJWT(token)
		if errParse != nil {
			http.Error(writer, custerrors.ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}
		ctxValue := context.WithValue(request.Context(), KeyIDUser, parseJwt.IdUser)
		ReqCtx := request.WithContext(ctxValue)
		next.ServeHTTP(writer, ReqCtx)
	})
}
func IsTempUser(next http.Handler, conf *configs.Config) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("X-User-Token")
		token, errValid := validateToken(authHeader)
		if errValid != nil {
			http.Error(writer, custerrors.ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}
		parseJwt, errParse := jwt.NewJWT(conf.Secret).ParseTemporaryJWT(token)
		if errParse != nil {
			http.Error(writer, custerrors.ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}
		ctxValue := context.WithValue(request.Context(), KeyTempUser, parseJwt.Email)
		ReqCtx := request.WithContext(ctxValue)
		next.ServeHTTP(writer, ReqCtx)
	})
}
