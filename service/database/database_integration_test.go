//go:build integration

package database_test

import (
	"context"
	"fmt"
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

func TestFind(t *testing.T) {
	ctx := context.Background()

	db := database.NewDatabase(config)
	err := db.Connect(ctx)
	if !assert.Nil(t, err) {
		return
	}

	data := database.Data{}

	dataFilter := database.DataFilter{
		ID:       true,
		Username: true,
		Password: true,
	}

	err = db.Find(ctx, &data, dataFilter)
	if !assert.Nil(t, err) {
		return
	}

	fmt.Println(data)

	err = db.Disconnect(ctx)
	if !assert.Nil(t, err) {
		return
	}
}
