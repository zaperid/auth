//go:build integration

package database_test

import (
	"context"
	"montrek-auth/service/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var config = database.Config{
	Logger:     zap.NewNop(),
	Host:       "mongodb://localhost:27017",
	Database:   "montrek",
	Collection: "user",
}

func TestConnection(t *testing.T) {
	ctx := context.Background()

	db := database.NewDatabase(config)
	err := db.Connect(ctx)
	if !assert.Nil(t, err) {
		return
	}

	err = db.Disconnect(ctx)
	if !assert.Nil(t, err) {
		return
	}
}
