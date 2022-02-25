package service

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type service_impl struct {
	logger *zap.Logger
	client *mongo.Client
	db     *mongo.Database
}

func NewService(logger *zap.Logger) Service {
	if logger == nil {
		logger = zap.NewNop()
	}

	service := service_impl{
		logger: logger,
	}
	return &service
}

func (service *service_impl) Connect(ctx context.Context, host string, database string) error {
	var err error

	service.logger.Info("connecting",
		zap.String("host", host),
		zap.String("database", database),
	)

	service.client, err = mongo.Connect(ctx, options.Client().ApplyURI(host))
	if err != nil {
		return err
	}

	service.db = service.client.Database(database)

	service.logger.Info("connected")
	return nil
}

func (service *service_impl) Disconnect(ctx context.Context) error {
	service.logger.Info("disconnecting")

	err := service.client.Disconnect(ctx)
	if err != nil {
		return err
	}

	service.logger.Info("disconnected")
	return nil
}
