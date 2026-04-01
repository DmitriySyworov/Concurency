package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapper := &WrapperWriter{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}
		next.ServeHTTP(wrapper, r)
		logger := logrus.New()
		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.WithFields(logrus.Fields{
			"leadTime":   time.Since(start).String(),
			"statusCode": wrapper.Status,
			"method":     r.Method,
			"path":       r.URL.Path,
		}).Info("Logging")
	})
}
