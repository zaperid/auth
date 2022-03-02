package rest

import "go.uber.org/zap"

type Server interface {
	Run() error
}

type Config struct {
	Logger *zap.Logger
	Port   uint16
}

type error_impl struct {
	Error string
}
