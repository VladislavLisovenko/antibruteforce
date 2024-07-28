package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/VladislavLisovenko/antibruteforce/internal/logger"
)

type HTTPWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (h *HTTPWriter) WriteHeader(statusCode int) {
	h.StatusCode = statusCode
	h.ResponseWriter.WriteHeader(statusCode)
}

func loggingMiddleware(_ context.Context, next http.Handler, logger logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		hw := &HTTPWriter{w, http.StatusOK}
		next.ServeHTTP(hw, r)

		logger.Info(
			fmt.Sprintf(
				"%s [%s] %s %s %s %d %f \"%s\"",
				r.RemoteAddr,
				time.Now(),
				r.Method,
				r.RequestURI,
				r.Proto,
				hw.StatusCode,
				time.Since(start).Seconds(),
				r.UserAgent(),
			),
		)
	})
}
