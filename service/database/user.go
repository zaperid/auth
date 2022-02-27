package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func (database *database_impl) Insert(ctx context.Context, data *Data) error {
	database.config.Logger.Info("insert data")
	database.config.Logger.Debug("inserting data", zap.Any("data", data))

	if database.col == nil {
		return ErrorDisconnected
	}

	result, err := database.col.InsertOne(ctx, data)
	if err != nil {
		switch err {
		case mongo.ErrClientDisconnected:
			return ErrorDisconnected
		default:
			database.config.Logger.Debug(ErrorNotIdentified.Error(), zap.String("error", err.Error()))
			return err
		}
	}

	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		data.ID = id
	}

	database.config.Logger.Debug("data inserted", zap.Any("data", data))
	return nil
}

func (database *database_impl) Find(ctx context.Context, data *Data) error {
	database.config.Logger.Info("find data")
	database.config.Logger.Debug("finding data", zap.Any("data", data))

	if database.col == nil {
		return ErrorDisconnected
	}

	result := database.col.FindOne(ctx, data)
	err := result.Decode(data)
	if err != nil {
		switch err {
		case mongo.ErrClientDisconnected:
			return ErrorDisconnected
		case mongo.ErrNoDocuments:
			return ErrorNotFound
		default:
			database.config.Logger.Debug(ErrorNotIdentified.Error(), zap.String("error", err.Error()))
			return err
		}
	}

	database.config.Logger.Debug("data found", zap.Any("data", data))
	return nil
}
