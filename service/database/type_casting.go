package database

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

var (
	ErrorNotIdentified = errors.New("error is not identifed")
	ErrorDisconnected  = errors.New("service disconnected")
	ErrorNotFound      = errors.New("data not found")
)

type Config struct {
	Logger     *zap.Logger
	Host       string
	Database   string
	Collection string
}

type Database interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Insert(ctx context.Context, data *Data) error
	Find(ctx context.Context, data *Data) error
}

type Data struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username,omitempty"`
	Password string             `bson:"password,omitempty"`
}
