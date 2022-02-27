package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

type jwt_impl struct {
	config Config
}

func NewJwt(config Config) JWT {
	jwt := jwt_impl{
		config: config,
	}

	return &jwt
}

func (jwtWrapper *jwt_impl) Generate(data Data) (string, error) {
	jwtWrapper.config.Logger.Info("generate token")
	jwtWrapper.config.Logger.Debug("generating jwt", zap.Any("data", data))

	claims := claims_impl{
		Data: data,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(jwtWrapper.config.Lifetime).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtWrapper.config.Key)
	if err != nil {
		jwtWrapper.config.Logger.Debug(ErrNotIdentified.Error(), zap.String("error", err.Error()))
		return "", ErrNotIdentified
	}

	jwtWrapper.config.Logger.Debug("generated", zap.String("token", tokenStr))
	return tokenStr, nil
}

func (jwtWrapper *jwt_impl) Verify(tokenStr string) bool {
	jwtWrapper.config.Logger.Info("verifying token")
	jwtWrapper.config.Logger.Debug("verifying", zap.String("token", tokenStr))
	token, err := jwt.ParseWithClaims(tokenStr, &claims_impl{}, func(token *jwt.Token) (interface{}, error) {
		return jwtWrapper.config.Key, nil
	})
	if err != nil {
		jwtWrapper.config.Logger.Debug("parse token error", zap.String("error", err.Error()))
		return false
	}

	if !token.Valid {
		jwtWrapper.config.Logger.Debug("token invalid")
		return false
	}

	claims, ok := token.Claims.(*claims_impl)
	if !ok {
		jwtWrapper.config.Logger.Debug("claims invalid")
		return false
	}

	if claims.ExpiresAt < time.Now().Unix() {
		jwtWrapper.config.Logger.Debug("token expired")
		return false
	}

	jwtWrapper.config.Logger.Debug("token valid")
	return true
}
