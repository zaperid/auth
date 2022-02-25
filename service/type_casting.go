package service

import "context"

type Service interface {
	Connect(ctx context.Context, host string, database string) error
	Disconnect(ctx context.Context) error
}
