package captcha

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

var (
	ErrNotIdentified = errors.New("error not identified")
)

type Captcha interface {
	Generate(height int, width int) (string, string, error)
	Verify(tokenStr string, answer string) bool
}

type Config struct {
	Logger     *zap.Logger
	Key        []byte
	Lifetime   time.Duration
	NoiseCount int
}

type claims_impl struct {
	Session    string `json:"session"`
	SessionKey string `json:"sesion_key"`
	jwt.StandardClaims
}
