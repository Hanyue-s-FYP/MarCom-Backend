package middleware

import "net/http"

// allow all by default
func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {    
        w.Header().Set("Access-Control-Allow-Origin", "*") // for now, for simplicity allow from all origin first
        w.Header().Add("Access-Control-Allow-Headers", "content-type")
		next.ServeHTTP(w, r)
	})
}
