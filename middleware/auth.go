package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/Hanyue-s-FYP/Marcom-Backend/utils"
)

// the routes that requires authentication
var authRoutes []string = []string{"/auth_test"}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if the route to be accessed requires authentication
		if slices.ContainsFunc(authRoutes, func(route string) bool {
			return strings.Contains(r.URL.Path, route)
		}) {
			// see if the Authorization header exists
			token := strings.Split(r.Header.Get("Authorization"), "Bearer ")

			if len(token) != 2 {
				utils.ResponseError(w, utils.HttpError{
					Code:       http.StatusUnauthorized,
					Message:    "Authentication token does not exist or is malformed",
					LogMessage: "no auth token",
				})
			} else {

				next.ServeHTTP(w, r)
			}
		}
	})
}
