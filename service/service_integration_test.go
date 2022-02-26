//go:build integration

package service_test

import (
	"context"
	"montrek-auth/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

const host = "mongodb://localhost:27017"
const database = "montrek"

func TestConnection(t *testing.T) {
	ctx := context.Background()

	service := service.NewService(nil)
	err := service.Connect(ctx, host, database)
	if !assert.Nil(t, err) {
		return
	}

	err = service.Disconnect(ctx)
	if !assert.Nil(t, err) {
		return
	}
}

func TestUser(t *testing.T) {
	ctx := context.Background()

	service := service.NewService(nil)
	user, err := service.User()
	if !assert.NotNil(t, err) {
		return
	}
	if !assert.Nil(t, user) {
		return
	}

	err = service.Connect(ctx, host, database)
	if !assert.Nil(t, err) {
		return
	}

	defer func() {
		err = service.Disconnect(ctx)
		if !assert.Nil(t, err) {
			return
		}
	}()

	user, err = service.User()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, user) {
		return
	}
}
