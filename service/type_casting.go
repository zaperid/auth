package service

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
)

var (
	ErrConnectToDatabase    = errors.New("unable connect to database")
	ErrDisconnectToDatabase = errors.New("unable disconnect form database")
	ErrGenerateCaptcha      = errors.New("unable generate captcha")
	ErrPassNotConfirm       = errors.New("password not confirm")
	ErrUsernamedUsed        = errors.New("username already been used")
	ErrCaptchaInvalid       = errors.New("captcha invalid")
	ErrLenUsername          = errors.New("username length must be between 4 to 20 characters")
	ErrLenPassword          = errors.New("password length must be between 8 to 20 characters")
	ErrInsertData           = errors.New("unable to insert data to database")
	ErrFindData             = errors.New("unable to find data from database")
	ErrPassword             = errors.New("password does not match")
	ErrGenerateToken        = errors.New("unable to generate jwt")
	ErrTokenInvalid         = errors.New("token invalid")
	ErrIDInvalid            = errors.New("error ID invalid")
	ErrOldPassword          = errors.New("old password does not match")
	ErrUpdateData           = errors.New("unable to update data to database")
)

type Service interface {
	Close() error
	GenerateCaptcha(height int, width int) (captchaToken string, image string, err error)
	UsedUsername(ctx context.Context, username string) (used bool, err error)
	Register(ctx context.Context, captchaToken string, captchaAnswer string, username string, password string, passwordConfirm string) (err error)
	Login(ctx context.Context, captchaToken string, captchaAnswer string, username string, password string) (token string, err error)
	ChangePassword(ctx context.Context, token string, captchaToken string, captchaAnswer string, currentPassword string, newPassword string, newPasswordConfirm string) (err error)
	RefreshToken(oldToken string) (token string, err error)
	UpdateProfile(ctx context.Context, token string, firstname string, lastname string, email string) (err error)
}

type Config struct {
	Logger            *zap.Logger
	DatabaseHost      string
	DatabaseName      string
	ColectionName     string
	Key               []byte
	CaptchaLifetime   time.Duration
	CaptchaNoiseCount int
	CaptchaShowline   int
	CaptchaLength     int
	JWTLifetime       time.Duration
}
