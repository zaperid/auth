//go:build unit

package jwt_test

import (
	"montrek-auth/service/jwt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestConsistency(t *testing.T) {
	data := jwt.Data{
		ID: "1234",
	}

	config := jwt.Config{
		Logger:   zap.NewNop(),
		Key:      []byte("4321"),
		Lifetime: 1 * time.Minute,
	}
	jwt := jwt.NewJwt(config)
	tokenStr, err := jwt.Generate(data)
	if !assert.Nil(t, err) {
		return
	}

	parsedData, valid := jwt.Parse(tokenStr)
	if !assert.True(t, valid) {
		return
	}

	if !assert.Equal(t, data, parsedData) {
		return
	}
}
