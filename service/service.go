package service

import (
	"context"
	"montrek-auth/service/captcha"
	"montrek-auth/service/database"
	"montrek-auth/service/hash"
	"montrek-auth/service/jwt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type service_impl struct {
	config  Config
	db      database.Database
	captcha captcha.Captcha
	jwt     jwt.JWT
	timeout time.Duration
}

func NewService(config Config) (Service, error) {
	defer config.Logger.Info("initalize service", zap.String("execution time", executionTime(time.Now())))

	dbConfig := database.Config{
		Logger:     config.Logger.Named("database"),
		Host:       config.DatabaseHost,
		Database:   config.DatabaseName,
		Collection: config.ColectionName,
	}

	captchaConfig := captcha.Config{
		Logger:     config.Logger.Named("captcha"),
		Key:        config.Key,
		Lifetime:   config.CaptchaLifetime,
		NoiseCount: config.CaptchaNoiseCount,
	}

	jwtConfig := jwt.Config{
		Logger:   config.Logger.Named("jwt"),
		Key:      config.Key,
		Lifetime: config.JWTLifetime,
	}

	service := service_impl{
		config:  config,
		db:      database.NewDatabase(dbConfig),
		captcha: captcha.NewCaptcha(captchaConfig),
		jwt:     jwt.NewJwt(jwtConfig),
		timeout: 30 * time.Second,
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
	defer service.config.Logger.Info("close service", zap.String("execution time", executionTime(time.Now())))

	ctx, ctxCancel := context.WithTimeout(context.Background(), service.timeout)
	err := service.db.Disconnect(ctx)
	ctxCancel()
	if err != nil {
		return err
	}

	return nil
}

func (service *service_impl) GenerateCaptcha(height int, width int) (token string, image string, err error) {
	defer service.config.Logger.Info("generate captcha", zap.String("execution time", executionTime(time.Now())))

	return service.captcha.Generate(height, width)
}

func (service *service_impl) Register(ctx context.Context, captchaToken string, answer string, username string, password string, passwordConfirm string) (err error) {
	defer service.config.Logger.Info("register user", zap.String("execution time", executionTime(time.Now())))

	if !service.captcha.Verify(captchaToken, answer) {
		return ErrCaptchaInvalid
	}

	{
		lenUsername := len(username)
		if lenUsername < 4 || lenUsername > 20 {
			return ErrLenUsername
		}
	}

	{
		lenPassword := len(password)
		if lenPassword < 8 || lenPassword > 20 {
			return ErrLenPassword
		}
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

func (service *service_impl) UsedUsername(ctx context.Context, username string) (used bool, err error) {
	defer service.config.Logger.Info("check username", zap.String("execution time", executionTime(time.Now())))

	data := database.Data{
		Username: username,
	}

	err = service.db.Find(ctx, &data)
	if err == database.ErrorNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (service *service_impl) Login(ctx context.Context, captchaToken string, answer string, username string, password string) (token string, err error) {
	defer service.config.Logger.Info("user login", zap.String("execution time", executionTime(time.Now())))

	if !service.captcha.Verify(captchaToken, answer) {
		return "", ErrCaptchaInvalid
	}

	data := database.Data{
		Username: username,
	}

	err = service.db.Find(ctx, &data)
	if err != nil {
		return "", err
	}

	if data.Password != hash.Hash([]byte(password)) {
		return "", ErrPassword
	}

	jwtData := jwt.Data{
		ID: data.ID.Hex(),
	}

	tokenStr, err := service.jwt.Generate(jwtData)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (service *service_impl) ChangePassword(ctx context.Context, token string, captchaToken string, captcha string, oldPassword string, newPassword string, newPasswordConfirm string) (err error) {
	defer service.config.Logger.Info("change user password", zap.String("execution time", executionTime(time.Now())))

	jwtData, valid := service.jwt.Parse(token)
	if !valid {
		return ErrTokenInvalid
	}

	if !service.captcha.Verify(captchaToken, captcha) {
		return ErrCaptchaInvalid
	}

	if newPassword != newPasswordConfirm {
		return ErrPassNotConfirm
	}

	var dbData database.Data

	dbData.ID, err = primitive.ObjectIDFromHex(jwtData.ID)
	if err != nil {
		return ErrIDInvalid
	}

	err = service.db.Find(ctx, &dbData)
	if err != nil {
		return err
	}

	if dbData.Password != hash.Hash([]byte(oldPassword)) {
		return ErrOldPassword
	}

	err = service.db.Update(ctx, database.Data{
		ID:       dbData.ID,
		Password: hash.Hash([]byte(newPassword)),
	})
	if err != nil {
		return err
	}

	return nil
}
