package user

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrorPassworNotMatch = errors.New("password does not match")
	ErrorNotFound        = errors.New("data not found")
)

type User interface {
	Register(ctx context.Context, data RegisterData) error
	Find(ctx context.Context, data FindData) (Data, error)
	UsedUsername(ctx context.Context, username string) (bool, error)
}

type RegisterData struct {
	Username        string
	Password        string
	ConfirmPassword string
}

type FindData struct {
	Username string
	Password string
}

type Data struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username,omitempty"`
	Password string             `bson:"password,omitempty"`
}
