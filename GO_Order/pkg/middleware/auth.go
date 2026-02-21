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

const KeyIDUser KeyCtx = "KeyIDUser"

func IsAuth(next http.Handler, conf *configs.Config) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(writer, custerrors.ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}
		if strings.Count(authHeader, ".") != 2 {
			http.Error(writer, custerrors.ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}
		parseJwt, errParse := jwt.NewJWT(conf.Secret).DecodeJWT(strings.TrimPrefix(authHeader, "Bearer "))
		if errParse != nil {
			http.Error(writer, custerrors.ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}
		ctxValue := context.WithValue(request.Context(), KeyIDUser, parseJwt.IdUser)
		ReqCtx := request.WithContext(ctxValue)
		next.ServeHTTP(writer, ReqCtx)
	})
}
