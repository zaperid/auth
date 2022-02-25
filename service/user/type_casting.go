package user

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrorPassworNotMatch = errors.New("password does not match")
)

type User interface {
	Register(ctx context.Context, data RegisterData) error
}

type RegisterData struct {
	Username        string
	Password        string
	ConfirmPassword string
}

type Data struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
}
