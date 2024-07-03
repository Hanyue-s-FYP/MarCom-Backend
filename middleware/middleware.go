package middleware

import "net/http"

type middleware func(http.Handler) http.Handler

// allows to chain multiple middlewares
func Use(ms ...middleware) middleware {
	return func(next http.Handler) http.Handler {
        // don't want to change the function parameter, store temp first
        // nextTemp := next
		for i := 0; i < len(ms); i++ {
            // record the current middleware
			m1 := ms[i]
            // apply the current middleware and update the middleware
            next = m1(next)
		}

        return next
	}
}
