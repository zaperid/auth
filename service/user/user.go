package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	used, err := user.UsedUsername(ctx, registerData.Username)
	if err != nil {
		return err
	}

	if used {
		return ErrorUsernameUsed
	}

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

func (user *user_impl) Find(ctx context.Context, findData FindData) (Data, error) {
	user.logger.Info("finding", zap.Any("data", findData))

	userData := Data{
		Username: findData.Username,
		Password: hash([]byte(findData.Password)),
	}

	result := user.collection.FindOne(ctx, userData)
	if result.Err() == mongo.ErrNoDocuments {
		return Data{}, ErrorNotFound
	}

	err := result.Decode(&userData)
	if err != nil {
		return Data{}, err
	}

	user.logger.Info("found", zap.Any("data", userData))

	return userData, nil
}

func (user *user_impl) UsedUsername(ctx context.Context, username string) (bool, error) {
	logger := user.logger.With(zap.String("username", username))
	logger.Info("checking username")

	var used bool

	data := Data{
		Username: username,
	}

	limit := int64(1)

	n, err := user.collection.CountDocuments(ctx, data, &options.CountOptions{Limit: &limit})
	if err != nil {
		return false, err
	}

	used = (n == 1)

	logger.Info("username checked", zap.Bool("used", used))
	return used, nil
}
