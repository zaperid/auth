package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User interface {
	Insert(ctx context.Context, data *Data) error
}

type Data struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
}
