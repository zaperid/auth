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

func (user *user_impl) Register(ctx context.Context, registerData RegisterData) error {
	user.logger.Info("register", zap.Any("data", registerData))

	if registerData.Password != registerData.ConfirmPassword {
		return ErrorPassworNotMatch
	}

	userData := Data{
		Username: registerData.Username,
		Password: hash([]byte(registerData.Password)),
	}

	result, err := user.collection.InsertOne(ctx, userData)
	if err != nil {
		return err
	}
	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		userData.ID = id
	}

	user.logger.Debug("registered", zap.Any("data", userData))
	return nil
}
