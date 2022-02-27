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
		Logger:   zap.NewExample(),
		Key:      []byte("4321"),
		Lifetime: 1 * time.Minute,
	}
	jwt := jwt.NewJwt(config)
	tokenStr, err := jwt.Generate(data)
	if !assert.Nil(t, err) {
		return
	}

	verifed := jwt.Verify(tokenStr)
	if !assert.True(t, verifed) {
		return
	}
}
