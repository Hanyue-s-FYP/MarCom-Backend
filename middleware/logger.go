package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// Middleware to log every request to the terminal
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s, Path: %s, Origin: %s\r\n", time.Now().Format("2006-01-02T15:04:05.000Z"), r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
