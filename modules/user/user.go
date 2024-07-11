// handles authentication (login, register, and user related stuff)
// route authentication is handled by the auth middleware
package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Hanyue-s-FYP/Marcom-Backend/utils"
	"github.com/golang-jwt/jwt/v5"
)

// stuffs embeded in the jwt token generated
type JWTClaims struct {
	UserID int
    Role int // TODO create the enum for roles
	jwt.RegisteredClaims
}

type LoginResponse struct {
	Token string
}

func Login(w http.ResponseWriter, r *http.Request) (*LoginResponse, error) {
	// generates jwt token with HS256 method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		UserID: 1, // TODO change to user id from db
        Role: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // default to 24 hour expiry
		},
	})

	if tokStr, err := token.SignedString([]byte("very-secure-key")); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "failed to generate jwt token",
			LogMessage: fmt.Sprintf("%v", err),
		}
	} else {
		return &LoginResponse{tokStr}, nil
	}
}
