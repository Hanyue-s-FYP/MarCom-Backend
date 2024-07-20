package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/Hanyue-s-FYP/Marcom-Backend/modules/user"
	"github.com/Hanyue-s-FYP/Marcom-Backend/utils"
	"github.com/golang-jwt/jwt/v5"
)

// the routes that requires authentication
var authRoutes []string = []string{
	"/auth_test", // just to test auth middleware is working, will remove
    "/get-me", // require token when getting themself
	"/products",   // everything related to product should be authenticated with business id (for now, revisit to allow for investors)
    "/agents", // same for agent
    "/environments", // same for environments
}

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
				jwtToken, err := jwt.ParseWithClaims(token[1], &user.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
					return []byte("very-secure-key"), nil // TODO change to get secret key
				})
				if err != nil {
					utils.ResponseError(w, utils.HttpError{
						Code:       http.StatusUnauthorized,
						Message:    "Authentication token does not exist or is malformed",
						LogMessage: err.Error(),
					})
				} else if claims, ok := jwtToken.Claims.(*user.JWTClaims); ok {
                    slog.Info(fmt.Sprintf("User ID: %d, Role: %d, Expires: %s", claims.UserID, claims.Role, claims.ExpiresAt))
					r.Header.Add("UserId", strconv.Itoa(claims.UserID))
                    r.Header.Add("Role", strconv.Itoa(int(claims.Role)))
					next.ServeHTTP(w, r)
				} else {
					utils.ResponseError(w, utils.HttpError{
						Code:       http.StatusUnauthorized,
						Message:    "Authentication token does not exist or is malformed",
						LogMessage: "authentication token does not exist or is malformed: unable to parse jwt claims",
					})
				}
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
