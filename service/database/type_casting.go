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
	Find(ctx context.Context, data *Data, filter DataFilter) error
	Update(ctx context.Context, data Data) error
}

type Data struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Username  string             `bson:"username,omitempty"`
	Password  string             `bson:"password,omitempty"`
	Firstname string             `bson:"firstname,omitempty"`
	Lastname  string             `bson:"lastname,omitempty"`
	Email     string             `bson:"email,omitempty"`
}

type DataFilter struct {
	ID        bool `bson:"_id,omitempty"`
	Username  bool `bson:"username,omitempty"`
	Password  bool `bson:"password,omitempty"`
	Firstname bool `bson:"firstname,omitempty"`
	Lastname  bool `bson:"lastname,omitempty"`
	Email     bool `bson:"email,omitempty"`
}
