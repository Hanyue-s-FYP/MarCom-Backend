package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
)

// Middleware to log every request to the terminal
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info(fmt.Sprintf("%s Path: %s, Origin: %s\r\n", r.Method, r.URL.Path, r.RemoteAddr))
		next.ServeHTTP(w, r)
	})
}
