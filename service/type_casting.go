package service

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

var (
	ErrPassNotConfirm = errors.New("password not confirm")
	ErrUsernamedUsed  = errors.New("username already been used")
)

type Service interface {
	Close() error
	Register(ctx context.Context, username string, password string, passwordConfirm string) error
}

type Config struct {
	Logger       *zap.Logger
	DatabaseHost string
	DatabaseName string
}
