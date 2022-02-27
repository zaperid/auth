package captcha

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

type captcha_impl struct {
	config Config
}

func NewCaptcha(config Config) Captcha {
	captch := captcha_impl{
		config: config,
	}

	return &captch
}

func (captcha *captcha_impl) Generate(height int, width int) (string, error) {
	captcha.config.Logger.Info("generating")

	driver := base64Captcha.NewDriverMath(height, width, captcha.config.NoiseCount, 0, nil, nil, nil)
	id, question, answer := driver.GenerateIdQuestionAnswer()
	questionImg, err := driver.DrawCaptcha(question)
	if err != nil {
		captcha.config.Logger.Debug(ErrNotIdentified.Error(), zap.String("error", err.Error()))
		return "", ErrNotIdentified
	}

	encrypted_answer := hash([]byte(id + answer))

	claims := claims_impl{
		Session:    id,
		SessionKey: encrypted_answer,
		Image:      questionImg.EncodeB64string(),
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(captcha.config.Lifetime).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(captcha.config.Key)
	if err != nil {
		captcha.config.Logger.Debug(ErrNotIdentified.Error(), zap.String("error", err.Error()))
		return "", ErrNotIdentified
	}

	captcha.config.Logger.Info("generated", zap.Any("clains", claims), zap.String("answer", answer))
	return tokenStr, nil
}

func (captcha *captcha_impl) Verify(tokenStr string, answer string) bool {
	logger := captcha.config.Logger.With(
		zap.String("signed string", tokenStr),
		zap.String("answer", answer),
	)

	logger.Info("verifying")
	token, err := jwt.ParseWithClaims(tokenStr, &claims_impl{}, func(token *jwt.Token) (interface{}, error) {
		return captcha.config.Key, nil
	})
	if err != nil {
		logger.Error("parse error")
		logger.Info("not verified")
		return false
	}

	if !token.Valid {
		logger.Error("token invalid")
		logger.Info("not verified")
		return false
	}

	claims, ok := token.Claims.(*claims_impl)
	if !ok {
		logger.Error("claims invalid")
		logger.Info("not verified")
		return false
	}

	if claims.ExpiresAt < time.Now().Unix() {
		logger.Error("claims expired")
		logger.Info("not verified")
		return false
	}

	encrypted_answer := hash([]byte(claims.Session + answer))
	if encrypted_answer != claims.SessionKey {
		logger.Error("answer does not match")
		logger.Info("not verified")
		return false
	}

	logger.Info("verified")
	return true
}
