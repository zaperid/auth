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
)

type Service interface {
	Close() error
	Register(ctx context.Context, token string, answer string, username string, password string, passwordConfirm string) error
	UsedUsername(ctx context.Context, username string) (bool, error)
	GenerateCaptcha(height int, width int) (string, error)
}

type Config struct {
	Logger            *zap.Logger
	DatabaseHost      string
	DatabaseName      string
	ColectionName     string
	Key               []byte
	CaptchaLifetime   time.Duration
	CaptchaNoiseCount int
}
