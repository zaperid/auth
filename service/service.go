package service

import (
	"context"
	"montrek-auth/service/captcha"
	"montrek-auth/service/database"
	"montrek-auth/service/hash"
	"time"
)

type service_impl struct {
	config  Config
	db      database.Database
	captcha captcha.Captcha
	timeout time.Duration
}

func NewService(config Config) (Service, error) {
	dbConfig := database.Config{
		Logger:     config.Logger,
		Host:       config.DatabaseHost,
		Database:   config.DatabaseName,
		Collection: config.ColectionName,
	}

	captchaConfig := captcha.Config{
		Logger:     config.Logger,
		Key:        config.Key,
		Lifetime:   config.CaptchaLifetime,
		NoiseCount: config.CaptchaNoiseCount,
	}

	service := service_impl{
		config:  config,
		db:      database.NewDatabase(dbConfig),
		timeout: 30 * time.Second,
		captcha: captcha.NewCaptcha(captchaConfig),
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), service.timeout)
	err := service.db.Connect(ctx)
	ctxCancel()
	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (service *service_impl) Close() error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), service.timeout)
	err := service.db.Disconnect(ctx)
	ctxCancel()
	if err != nil {
		return err
	}

	return nil
}

func (service *service_impl) GenerateCaptcha(height int, width int) (string, error) {
	return service.captcha.Generate(height, width)
}

func (service *service_impl) Register(ctx context.Context, captchaToken string, answer string, username string, password string, passwordConfirm string) error {
	if !service.captcha.Verify(captchaToken, answer) {
		return ErrCaptchaInvalid
	}

	if password != passwordConfirm {
		return ErrPassNotConfirm
	}

	used, err := service.UsedUsername(ctx, username)
	if err != nil {
		return err
	}
	if used {
		return ErrUsernamedUsed
	}

	data := database.Data{
		Username: username,
		Password: hash.Hash([]byte(password)),
	}

	err = service.db.Insert(ctx, &data)
	if err != nil {
		return err
	}

	return nil
}

func (service *service_impl) UsedUsername(ctx context.Context, username string) (bool, error) {
	data := database.Data{
		Username: username,
	}

	err := service.db.Find(ctx, &data)
	if err != database.ErrorNotFound {
		return true, nil
	}
	if err != nil {
		return false, err
	}

	return false, nil
}
