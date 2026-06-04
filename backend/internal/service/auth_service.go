package service

import (
	"errors"
	"strconv"
	"unicode"

	"cowork/internal/config"
	"cowork/internal/model"
	"cowork/internal/repository"
	"cowork/pkg/captcha"
	"cowork/pkg/errcode"
	"cowork/pkg/jwt"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService handles authentication business logic.
type AuthService struct {
	DB    *gorm.DB
	Redis *redis.Client
	JWT   *config.JWTConfig
}

// LoginResp represents the login/refresh response data.
type LoginResp struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         *model.User `json:"user"`
}

// Register creates a new user account.
func (s *AuthService) Register(req RegisterReq) (*model.User, error) {
	// Check username uniqueness
	existing, err := repository.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, newAppError(errcode.ErrUserExists, "username already exists")
	}

	// Check email uniqueness
	existing, err = repository.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, newAppError(errcode.ErrUserExists, "email already exists")
	}

	// Validate password strength
	if !isStrongPassword(req.Password) {
		return nil, newAppError(errcode.ErrWeakPassword, "password must be at least 8 characters with uppercase, lowercase, and digit")
	}

	// Hash password with bcrypt cost 12
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedBytes),
		Role:         "Developer",
	}

	if err := repository.CreateUser(user); err != nil {
		return nil, err
	}

	// Clear password hash before returning
	user.PasswordHash = ""
	return user, nil
}

// Login authenticates a user and returns tokens.
func (s *AuthService) Login(req LoginReq, captchaID, captchaAnswer string) (*LoginResp, error) {
	// Verify captcha
	if !captcha.VerifyCaptcha(s.Redis, captchaID, captchaAnswer) {
		return nil, newAppError(errcode.ErrCaptchaInvalid, "invalid captcha")
	}

	// Find user by username
	user, err := repository.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, newAppError(errcode.ErrUserNotFound, "user not found")
	}

	// Compare bcrypt password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, newAppError(errcode.ErrWrongPassword, "wrong password")
	}

	// Generate tokens
	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Username, user.Role, s.JWT.AccessSecret, s.JWT.AccessTTL)
	if err != nil {
		return nil, err
	}
	refreshToken, err := jwt.GenerateRefreshToken(user.ID, s.JWT.RefreshSecret, s.JWT.RefreshTTL)
	if err != nil {
		return nil, err
	}

	// Clear password hash before returning
	user.PasswordHash = ""
	return &LoginResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// RefreshToken parses a refresh token and issues a new token pair.
func (s *AuthService) RefreshToken(refreshToken string) (*LoginResp, error) {
	// Parse refresh token
	claims, err := jwt.ParseToken(refreshToken, s.JWT.RefreshSecret)
	if err != nil {
		return nil, newAppError(errcode.ErrTokenInvalid, "invalid refresh token")
	}

	// Find user by ID from claims (refresh token stores userID in Subject)
	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return nil, newAppError(errcode.ErrTokenInvalid, "invalid user ID in token")
	}
	user, err := repository.FindByID(uint(userID))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, newAppError(errcode.ErrUserNotFound, "user not found")
	}

	// Generate new token pair
	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Username, user.Role, s.JWT.AccessSecret, s.JWT.AccessTTL)
	if err != nil {
		return nil, err
	}
	newRefreshToken, err := jwt.GenerateRefreshToken(user.ID, s.JWT.RefreshSecret, s.JWT.RefreshTTL)
	if err != nil {
		return nil, err
	}

	// Clear password hash before returning
	user.PasswordHash = ""
	return &LoginResp{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}

// GetProfile retrieves a user profile by ID.
func (s *AuthService) GetProfile(userID uint) (*model.User, error) {
	user, err := repository.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, newAppError(errcode.ErrUserNotFound, "user not found")
	}

	// Clear password hash before returning
	user.PasswordHash = ""
	return user, nil
}

// GenerateAndStoreCaptcha generates a captcha and stores the answer in Redis.
func (s *AuthService) GenerateAndStoreCaptcha() (id, b64 string, err error) {
	id, b64, answer, err := captcha.GenerateCaptcha()
	if err != nil {
		return "", "", err
	}
	if err := captcha.StoreCaptcha(s.Redis, id, answer); err != nil {
		return "", "", err
	}
	return id, b64, nil
}

// RegisterReq is the request for registration (used by handler, defined here for convenience).
type RegisterReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginReq is the request for login (used internally by service).
type LoginReq struct {
	Username      string
	Password      string
	CaptchaID     string
	CaptchaAnswer string
}

// isStrongPassword checks that password is >=8 chars, has uppercase, lowercase, and digit.
func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	var hasUpper, hasLower, hasDigit bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	return hasUpper && hasLower && hasDigit
}

// AppError is a custom error type that carries an error code.
type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func newAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// IsAppError checks if an error is an AppError and returns it.
func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
