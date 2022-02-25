//go:build integration

package service_test

import (
	"context"
	"montrek-api/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnection(t *testing.T) {
	const host = "mongodb://localhost:27017"
	const database = "montrek"

	ctx := context.Background()

	service := service.NewService(nil)
	err := service.Connect(ctx, host, database)
	if assert.Nil(t, err) {
		return
	}

	err = service.Disconnect(ctx)
	if assert.Nil(t, err) {
		return
	}
}
