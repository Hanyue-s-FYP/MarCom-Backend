// handles authentication (login, register, and user related stuff)
// route authentication is handled by the auth middleware
package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Hanyue-s-FYP/Marcom-Backend/db"
	"github.com/Hanyue-s-FYP/Marcom-Backend/db/models"
	"github.com/Hanyue-s-FYP/Marcom-Backend/modules"
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

// allows posting to this route to create a business account, validations are done at front end
func RegisterBusiness(w http.ResponseWriter, r *http.Request) (*modules.ExecResponse, error) {
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

	return &modules.ExecResponse{Message: "Successfully registered a business account, please login"}, nil
}

type LoginUserForm struct {
	Username string
	Password string
}

type LoginResponse struct {
	Token   string
	Message string
}

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
				Code:       http.StatusUnauthorized,
				Message:    "Failed to login, please check credentials",
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
			Code:       http.StatusUnauthorized,
			Message:    "Failed to login, please check credentials",
			LogMessage: fmt.Sprintf("failed to login: %v", err),
		}
	}

	// generates jwt token with HS256 method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		UserID: user.ID,         // depending on role if business then is BusinessID
		Role:   models.BUSINESS, // default to business first now dont care about investor gok
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // default to 24 hour expiry
		},
	})

	config := utils.GetConfig()
	if tokStr, err := token.SignedString([]byte(config.JWT_SECRET_KEY)); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "failed to generate jwt token",
			LogMessage: fmt.Sprintf("%v", err),
		}
	} else {
		return &LoginResponse{Token: tokStr, Message: "Successfully logged in"}, nil
	}
}

type UserWithoutPassword struct {
	ID          int
	Username    string
	DisplayName string
	Email       string
	Status      int
	PhoneNumber string
}

type UserWithRole struct {
	User   UserWithoutPassword
	Role   models.UserRole
	RoleID int //denotes the id for that specific role, if role is business then is business id vice versa
}

func GetMe(w http.ResponseWriter, r *http.Request) (*UserWithRole, error) {
	// obtain role from header (passed through the auth middleware alrd)
	role, err := strconv.Atoi(r.Header.Get("Role"))
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain user account",
			LogMessage: fmt.Sprintf("failed to obtain user role when get user: %v", err),
		}
	}

	var id, userId int
	if id, err = strconv.Atoi(r.Header.Get("UserId")); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain user account",
			LogMessage: fmt.Sprintf("failed to obtain user id when get user: %v", err),
		}
	}
	if role == models.INVESTOR {
		// do fetch user through investor model
	} else {
		business, err := models.BusinessModel.GetByBusinessID(id)
		if err != nil {
			return nil, utils.HttpError{
				Code:       http.StatusInternalServerError,
				Message:    "Failed to obtain user account",
				LogMessage: fmt.Sprintf("failed to obtain business when get user: %v", err),
			}
		}
		userId = business.User.ID
	}

	user, err := models.UserModel.GetByID(userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, utils.HttpError{
				Code:       http.StatusNotFound,
				Message:    "User does not exist",
				LogMessage: "user not found",
			}
		}
	}

	return &UserWithRole{
		User: UserWithoutPassword{
			ID:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Email:       user.Email,
			Status:      user.Status,
			PhoneNumber: user.PhoneNumber,
		},
		Role:   models.UserRole(role),
		RoleID: id,
	}, nil
}

type BusinessWithoutPassword struct {
	User         UserWithoutPassword
	ID           int
	Description  string
	BusinessType string
	CoverImgPath string
}

func GetBusiness(w http.ResponseWriter, r *http.Request) (*BusinessWithoutPassword, error) {
	idStr := r.PathValue("id") // assumes {id} exists in the route
	if len(idStr) == 0 {
		return nil, utils.HttpError{
			Code:       http.StatusBadRequest,
			Message:    "Expected business id in path, got none",
			LogMessage: "got empty path value when obtaining business",
		}
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusBadRequest,
			Message:    "Failed to parse business id to number",
			LogMessage: fmt.Sprintf("failed to parse business id from path value: %v", err),
		}
	}

	business, err := models.BusinessModel.GetByBusinessID(id)
	if err != nil {
		if errors.Is(err, models.ErrBusinessNotFound) {
			return nil, utils.HttpError{
				Code:       http.StatusNotFound,
				Message:    "Failed to obtain business, business does not exist",
				LogMessage: fmt.Sprintf("failed to obtain business: %v", err),
			}
		}
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to obtain business",
			LogMessage: fmt.Sprintf("failed to obtain business: %v", err),
		}
	}
	return &BusinessWithoutPassword{
		User: UserWithoutPassword{
			ID:          business.User.ID,
			Username:    business.User.Username,
			DisplayName: business.User.DisplayName,
			Email:       business.User.Email,
			Status:      business.User.Status,
			PhoneNumber: business.User.PhoneNumber,
		},
		ID:           business.ID,
		Description:  business.Description,
		BusinessType: business.BusinessType,
		CoverImgPath: business.CoverImgPath,
	}, nil
}

func CheckUserWithUsername(w http.ResponseWriter, r *http.Request) (*UserWithoutPassword, error) {
	username := r.PathValue("username")

	// check if account already exist (username must be unique)
	if user, err := models.BusinessModel.GetByUsername(username); user != nil {
		return &UserWithoutPassword{
			ID:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Email:       user.Email,
			Status:      user.Status,
			PhoneNumber: user.PhoneNumber,
		}, nil
	} else if err != nil {
		if errors.Is(err, models.ErrBusinessNotFound) {
			return nil, nil // front end check error, no want use error
		}
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Something unexpected happened when checking username uniqueness",
			LogMessage: fmt.Sprintf("something failed when checking username uniqueness: %v", err),
		}
	}

	// not very possible reach this point
	return nil, nil
}

func UpdateBusiness(w http.ResponseWriter, r *http.Request) (*modules.ExecResponse, error) {
	var business models.Business
	r.ParseMultipartForm(1 << 30) // 1GB max size should be sufficient

	idStr := r.FormValue("ID")
	slog.Info(fmt.Sprintf("Updating business with ID: %s", idStr))
	if id, err := strconv.Atoi(idStr); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to parse business ID",
			LogMessage: fmt.Sprintf("failed to parse business ID when updating business: %v", err),
		}
	} else {
		business.ID = id
	}

	// get all the other form values that can be updated
	business.Description = r.FormValue("Description")
	business.CoverImgPath = r.FormValue("CoverImgPath")

	// front end will handle new cover image will send back through another property for easy purpose (no need mess with complex typing)
	// the CoverImgPath property should be left unchanged and remain original when sent back from front end
	file, header, err := r.FormFile("NewCoverImg")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Unexpected error occured",
			LogMessage: fmt.Sprintf("unexpected error when updating business: %v", err),
		}
	}

	// only if got file, handle upload file
	if !errors.Is(err, http.ErrMissingFile) {
		slog.Info(fmt.Sprintf("Obtained file, Filename: %s", header.Filename))
		uploadPath, err := db.UploadImage(&file, header)
		if err != nil {
			return nil, utils.HttpError{
				Code:       http.StatusInternalServerError,
				Message:    "Failed to upload image",
				LogMessage: fmt.Sprintf("failed to upload file to database: %v", err),
			}
		}
		slog.Info(fmt.Sprintf("File uploaded to: %s", uploadPath))
		business.CoverImgPath = uploadPath
	}

	// only should be cannot get count
	if _, err := models.BusinessModel.GetByBusinessID(business.ID); err != nil {
		if errors.Is(err, models.ErrBusinessNotFound) {
			return nil, utils.HttpError{
				Code:       http.StatusNotFound,
				Message:    "Business to update does not exist",
				LogMessage: fmt.Sprintf("business with id %d not found", business.ID),
			}
		} else {
			return nil, utils.HttpError{
				Code:       http.StatusInternalServerError,
				Message:    "Unexpected error occured",
				LogMessage: fmt.Sprintf("unexpected error when updating business: %v", err),
			}
		}
	}
	if err := models.BusinessModel.Update(business); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to update business",
			LogMessage: fmt.Sprintf("failed to update business: %v", err),
		}
	}

	return &modules.ExecResponse{Message: "Successfully updated business"}, nil
}

type ForgetPasswordData struct {
	Username string
	Email    string
}

type ForgetPasswordClaims struct {
	Username string
	UserID   int
	jwt.RegisteredClaims
}

func ForgetPassword(w http.ResponseWriter, r *http.Request) (*modules.ExecResponse, error) {
	var forgetPasswordData ForgetPasswordData
	if err := json.NewDecoder(r.Body).Decode(&forgetPasswordData); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusBadRequest,
			Message:    "Failed to parse username and email from request",
			LogMessage: fmt.Sprintf("failed to decode request for forget password data: %v", err),
		}
	}

	user, err := models.UserModel.GetByUsername(forgetPasswordData.Username)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, utils.HttpError{
				Code:    http.StatusNotFound,
				Message: "Username not found in the system",
			}
		} else {
			return nil, utils.HttpError{
				Code:       http.StatusInternalServerError,
				Message:    "Failed to obtain user based on username",
				LogMessage: fmt.Sprintf("failed to obtain user by username: %v", err),
			}
		}
	}

	// check if email match
	if user.Email != forgetPasswordData.Email {
		return nil, utils.HttpError{
			Code:    http.StatusForbidden,
			Message: "User is not registered under this email",
		}
	}

	// generate jwt and send email
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, ForgetPasswordClaims{
		Username: user.Username,
		UserID:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // token expires in 15 minutes
		},
	})

	config := utils.GetConfig()
	tokStr, err := token.SignedString([]byte(config.JWT_SECRET_KEY))
	if err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to generate jwt token",
			LogMessage: fmt.Sprintf("%v", err),
		}
	}

	if err = utils.SendMail(user.Email, "Reset Password", fmt.Sprintf("Hi %s, please reset your password through this link: %s/reset-password/%s\nPlease note that the link will expire in 15 minutes. If you did not make this request, you may safely ignore this email.", user.DisplayName, config.FRONT_END_ADDR, tokStr)); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to send email",
			LogMessage: fmt.Sprintf("failed to send email: %v", err),
		}
	}
	return &modules.ExecResponse{Message: fmt.Sprintf("Successfully sent reset password link to %s", user.Email)}, nil
}

type ResetForgetPassword struct {
	Password            string
	ForgetPasswordToken string
}

func ResetPassword(w http.ResponseWriter, r *http.Request) (*modules.ExecResponse, error) {
	var resetForgetPassword ResetForgetPassword

	if err := json.NewDecoder(r.Body).Decode(&resetForgetPassword); err != nil {
		return nil, utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to parse request body",
			LogMessage: fmt.Sprintf("failed to parse request body: %v", err),
		}
	}

	config := utils.GetConfig()
	jwtToken, err := jwt.ParseWithClaims(resetForgetPassword.ForgetPasswordToken, &ForgetPasswordClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JWT_SECRET_KEY), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, utils.HttpError{
				Code:       http.StatusUnauthorized,
				Message:    "Reset password token expired, please login again",
				LogMessage: err.Error(),
			}
		} else {
			return nil, utils.HttpError{
				Code:       http.StatusUnauthorized,
				Message:    "Reset password token does not exist or is malformed",
				LogMessage: err.Error(),
			}
		}
	} else if claims, ok := jwtToken.Claims.(*ForgetPasswordClaims); ok {
		if err := changePassword(claims.UserID, resetForgetPassword.Password); err != nil {
			return nil, err
		}

		return &modules.ExecResponse{Message: "Successfully reset password"}, nil
	} else {
		return nil, utils.HttpError{
			Code:       http.StatusUnauthorized,
			Message:    "Reset password token does not exist or is malformed",
			LogMessage: "reset password token does not exist or is malformed: unable to parse jwt claims",
		}
	}
}

func hashPassword(pw string) (string, error) {
	if hashed, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost); err != nil {
		return "", err
	} else {
		return string(hashed), nil
	}
}

func changePassword(userId int, pw string) error {
	user, err := models.UserModel.GetByID(userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return utils.HttpError{
				Code:    http.StatusNotFound,
				Message: "Username not found in the system",
			}
		} else {
			return utils.HttpError{
				Code:       http.StatusInternalServerError,
				Message:    "Failed to obtain user based on user ID",
				LogMessage: fmt.Sprintf("failed to obtain user by user ID: %v", err),
			}
		}
	}

	hashedPw, err := hashPassword(pw)
	if err != nil {
		return utils.HttpError{
			Code:       http.StatusInternalServerError,
			Message:    "Failed to create a business account",
			LogMessage: fmt.Sprintf("failed to hash password: %v", err),
		}
	}
	user.Password = hashedPw

	if err = models.UserModel.Update(*user); err != nil {
		return err
	}

	return nil
}
