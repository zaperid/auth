package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

var (
	ErrNotIdentified = errors.New("error not identified")
)

type JWT interface {
	Generate(data Data) (token string, err error)
	Parse(token string) (data Data, valid bool)
}

type Config struct {
	Logger   *zap.Logger
	Key      []byte
	Lifetime time.Duration
}

type Data struct {
	ID string
}

type claims_impl struct {
	Data
	jwt.RegisteredClaims
}
