package captcha

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

type captcha_impl struct {
	config Config
	source string
}

func NewCaptcha(config Config) Captcha {
	captch := captcha_impl{
		config: config,
		source: "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
	}

	return &captch
}

func (captcha *captcha_impl) Generate(height int, width int) (tokenStr string, imgStr string, err error) {
	captcha.config.Logger.Info("generate captcha")
	captcha.config.Logger.Debug("generating captcha", zap.Int("height", height), zap.Int("width", width))

	driver := base64Captcha.NewDriverString(
		height, width,
		captcha.config.NoiseCount,
		captcha.config.ShowLine,
		captcha.config.Length,
		captcha.source,
		nil, nil, nil,
	)
	id, question, answer := driver.GenerateIdQuestionAnswer()
	img, err := driver.DrawCaptcha(question)
	if err != nil {
		captcha.config.Logger.Debug(ErrNotIdentified.Error(), zap.String("error", err.Error()))
		return "", "", ErrNotIdentified
	}
	imgStr = img.EncodeB64string()

	encrypted_answer := hash([]byte(id + answer))

	claims := claims_impl{
		Session:    id,
		SessionKey: encrypted_answer,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(captcha.config.Lifetime)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = token.SignedString(captcha.config.Key)
	if err != nil {
		captcha.config.Logger.Debug(ErrNotIdentified.Error(), zap.String("error", err.Error()))
		return "", "", ErrNotIdentified
	}

	captcha.config.Logger.Debug("captcha generated", zap.Any("token", tokenStr), zap.String("image", imgStr), zap.String("answer", answer))
	return tokenStr, imgStr, nil
}

func (captcha *captcha_impl) Verify(tokenStr string, answer string) (valid bool) {
	captcha.config.Logger.Info("verify captcha")
	captcha.config.Logger.Debug("verify captcha",
		zap.String("token", tokenStr),
		zap.String("answer", answer),
	)

	token, err := jwt.ParseWithClaims(tokenStr, &claims_impl{}, func(token *jwt.Token) (interface{}, error) {
		return captcha.config.Key, nil
	})
	if err != nil {
		switch err {
		default:
			captcha.config.Logger.Debug("unidentified parse captcha's token error", zap.String("error", err.Error()))
			return false
		}
	}

	if !token.Valid {
		captcha.config.Logger.Debug("captcha's token invalid")
		return false
	}

	claims, ok := token.Claims.(*claims_impl)
	if !ok {
		captcha.config.Logger.Debug("captcha's claims invalid")
		return false
	}

	encrypted_answer := hash([]byte(claims.Session + answer))
	if encrypted_answer != claims.SessionKey {
		captcha.config.Logger.Debug("wrong captcha answer")
		return false
	}

	captcha.config.Logger.Debug("captcha valid")
	return true
}
