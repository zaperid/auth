package service

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
)

var (
	ErrPassNotConfirm = errors.New("password not confirm")
	ErrUsernamedUsed  = errors.New("username already been used")
	ErrCaptchaInvalid = errors.New("captcha invalid")
	ErrLenUsername    = errors.New("username length must be between 4 to 20 characters")
	ErrLenPassword    = errors.New("password length must be between 8 to 20 characters")
)

type Service interface {
	Close() error
	GenerateCaptcha(height int, width int) (token string, image string, err error)
	Register(ctx context.Context, captchaToken string, answer string, username string, password string, passwordConfirm string) (err error)
	UsedUsername(ctx context.Context, username string) (used bool, err error)
	Login(ctx context.Context, captchaToken string, answer string, username string, password string) (token string, err error)
}

type Config struct {
	Logger            *zap.Logger
	DatabaseHost      string
	DatabaseName      string
	ColectionName     string
	Key               []byte
	CaptchaLifetime   time.Duration
	CaptchaNoiseCount int
	JWTLifetime       time.Duration
}
