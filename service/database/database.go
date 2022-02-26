package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type database_impl struct {
	config Config
	client *mongo.Client
	col    *mongo.Collection
}

func NewDatabase(config Config) Database {
	database := database_impl{
		config: config,
	}
	return &database
}

func (database *database_impl) Connect(ctx context.Context) error {
	var err error

	database.config.Logger.Info("connecting")

	database.client, err = mongo.Connect(ctx, options.Client().ApplyURI(database.config.Host))
	if err != nil {
		database.config.Logger.Debug(ErrorNotIdentified.Error(), zap.String("error", err.Error()))
		return ErrorNotIdentified
	}

	database.col = database.client.Database(database.config.Database).Collection(database.config.Collection)

	database.config.Logger.Info("connected")
	return nil
}

func (database *database_impl) Disconnect(ctx context.Context) error {
	database.config.Logger.Info("disconnecting")

	err := database.client.Disconnect(ctx)
	if err != nil {
		database.config.Logger.Debug(ErrorNotIdentified.Error(), zap.String("error", err.Error()))
		return ErrorNotIdentified
	}

	database.col = nil

	database.config.Logger.Info("disconnected")
	return nil
}
