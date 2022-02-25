package service

import (
	"context"
	"errors"
	"montrek-api/service/user"
)

var (
	ErrorDisconnected = errors.New("service disconnected")
)

type Service interface {
	Connect(ctx context.Context, host string, database string) error
	Disconnect(ctx context.Context) error
	User() (user.User, error)
}
