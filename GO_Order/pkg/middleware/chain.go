package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for _, mv := range middlewares{
			next = mv(next)
		}
		return next
	}
}
