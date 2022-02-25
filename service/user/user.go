package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type user_impl struct {
	logger     *zap.Logger
	collection *mongo.Collection
}

func NewUser(logger *zap.Logger, collection *mongo.Collection) User {
	user := user_impl{
		logger:     logger,
		collection: collection,
	}

	return &user
}

func (user *user_impl) Insert(ctx context.Context, data *Data) error {
	user.logger.Info("inserting", zap.Any("data", data))

	data.Password = hash([]byte(data.Password))

	result, err := user.collection.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		data.ID = id
	}

	user.logger.Debug("inserted", zap.Any("data", data))

	return nil
}
