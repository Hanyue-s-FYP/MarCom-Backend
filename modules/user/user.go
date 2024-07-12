// handles authentication (login, register, and user related stuff)
// route authentication is handled by the auth middleware
package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Hanyue-s-FYP/Marcom-Backend/db/models"
	"github.com/Hanyue-s-FYP/Marcom-Backend/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// stuffs embeded in the jwt token generated
type JWTClaims struct {
	UserID int
	Role   models.UserRole
	jwt.RegisteredClaims
}

type RegisterResponse struct {
	Message string
}

// allows posting to this route to create a business account, validations are done at front end
// TODO revisit when have more time and do backend validations as well
func RegisterBusiness(w http.ResponseWriter, r *http.Request) (*RegisterResponse, error) {
	var user models.Business

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to create a business account",
			LogMessage: fmt.Sprintf("failed to decode business account body: %v", err),
		}
	}

	// check if account already exist (username must be unique)
	if business, err := models.BusinessModel.GetByUsername(user.Username); business != nil {
		return nil, utils.HttpError{
			Code:    http.StatusConflict,
			Message: fmt.Sprintf("User with username %s already exist", user.Username),
			LogMessage: utils.If(
				err != nil,
				fmt.Sprintf("failed to create business account: %v", err),
				"failed to create business account: username already exist",
			), // not necessary is account exist, but mask that away from user lest they get frustrated
		}
	}

	// hash password using bcrypt
	hashedPw, err := hashPassword(user.Password)
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to create a business account",
			LogMessage: fmt.Sprintf("failed to hash password: %v", err),
		}
	}
	user.Password = hashedPw

	if err := models.BusinessModel.Create(user); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to create a business account",
			LogMessage: fmt.Sprintf("failed to create business account: %v", err),
		}
	}

	return &RegisterResponse{"Successfully registered a business account, please login"}, nil
}

type LoginUserForm struct {
	Username string
	Password string
}

type LoginResponse struct {
	Token   string
	Message string
}

// TODO handle multiple roles (xian complete business one ka lai do investor eh d)
func Login(w http.ResponseWriter, r *http.Request) (*LoginResponse, error) {
	// take body posted and check if any data in sqlite
	var creds LoginUserForm
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to parse login data",
			LogMessage: fmt.Sprintf("failed to parse login data: %v", err),
		}
	}

	user, err := models.BusinessModel.GetByUsername(creds.Username)
	if err != nil {
		if errors.Is(err, models.ErrBusinessNotFound) {
			return nil, utils.HttpError{
				Code:    http.StatusUnauthorized,
				Message: "Failed to login, please check credentials",
                LogMessage: "failed to login, account not found",
			}
		}
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to login, something unexpected happened, please wait while we try to fix it",
			LogMessage: fmt.Sprintf("failed to login: %v", err),
		}
	}

	// verify if the password can match stored hash, bcrypt does not recommend hashing the password again, instead, use ComparePasswordWithHash (gonna read more on this)
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		return nil, utils.HttpError{
			Code:    http.StatusUnauthorized,
			Message: "Failed to login, please check credentials",
            LogMessage: fmt.Sprintf("failed to login: %v", err),
		}
	}

	// generates jwt token with HS256 method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		UserID: 1,
		Role:   1,
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
		return &LoginResponse{Token: tokStr, Message: "Successfully logged in"}, nil
	}
}

func hashPassword(pw string) (string, error) {
	if hashed, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost); err != nil {
		return "", err
	} else {
		return string(hashed), nil
	}
}
