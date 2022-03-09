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
		ShowLine:   config.CaptchaShowline,
		Length:     config.CaptchaLength,
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
		service.config.Logger.Error(err.Error())
		return nil, ErrConnectToDatabase
	}

	return &service, nil
}

func (service *service_impl) Close() error {
	defer service.config.Logger.Info("close service", zap.String("execution time", executionTime(time.Now())))

	ctx, ctxCancel := context.WithTimeout(context.Background(), service.timeout)
	err := service.db.Disconnect(ctx)
	ctxCancel()

	if err != nil {
		service.config.Logger.Error(err.Error())
		return ErrConnectToDatabase
	}

	return nil
}

func (service *service_impl) GenerateCaptcha(height int, width int) (captchaToken string, image string, err error) {
	defer service.config.Logger.Info("generate captcha", zap.String("execution time", executionTime(time.Now())))

	captchaToken, image, err = service.captcha.Generate(height, width)
	if err != nil {
		return "", "", ErrGenerateCaptcha
	}

	return captchaToken, image, nil
}

func (service *service_impl) Register(ctx context.Context, captchaToken string, captchaAnswer string, username string, password string, passwordConfirm string) (err error) {
	defer service.config.Logger.Info("register user", zap.String("execution time", executionTime(time.Now())))

	if !service.captcha.Verify(captchaToken, captchaAnswer) {
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
		service.config.Logger.Error(err.Error())
		return ErrInsertData
	}

	return nil
}

func (service *service_impl) UsedUsername(ctx context.Context, username string) (used bool, err error) {
	defer service.config.Logger.Info("check username", zap.String("execution time", executionTime(time.Now())))

	data := database.Data{
		Username: username,
	}

	filter := database.DataFilter{
		ID: true,
	}

	err = service.db.Find(ctx, &data, filter)
	if err == database.ErrorNotFound {
		return false, nil
	}
	if err != nil {
		service.config.Logger.Error(err.Error())
		return false, ErrFindData
	}

	return true, nil
}

func (service *service_impl) Login(ctx context.Context, captchaToken string, captchaAnswer string, username string, password string) (token string, err error) {
	defer service.config.Logger.Info("user login", zap.String("execution time", executionTime(time.Now())))

	if !service.captcha.Verify(captchaToken, captchaAnswer) {
		return "", ErrCaptchaInvalid
	}

	data := database.Data{
		Username: username,
	}

	filter := database.DataFilter{
		ID:       true,
		Password: true,
	}

	err = service.db.Find(ctx, &data, filter)
	if err != nil {
		service.config.Logger.Error(err.Error())
		return "", ErrFindData
	}

	if data.Password != hash.Hash([]byte(password)) {
		return "", ErrPassword
	}

	jwtData := jwt.Data{
		ID: data.ID.Hex(),
	}

	tokenStr, err := service.jwt.Generate(jwtData)
	if err != nil {
		service.config.Logger.Error(err.Error())
		return "", ErrGenerateToken
	}

	return tokenStr, nil
}

func (service *service_impl) ChangePassword(ctx context.Context, token string, captchaToken string, captcha string, currentPassword string, newPassword string, newPasswordConfirm string) (err error) {
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

	var data database.Data
	{
		data.ID, err = primitive.ObjectIDFromHex(jwtData.ID)
		if err != nil {
			service.config.Logger.Error(err.Error())
			return ErrIDInvalid
		}
	}

	filter := database.DataFilter{
		ID:       true,
		Password: true,
	}

	err = service.db.Find(ctx, &data, filter)
	if err != nil {
		service.config.Logger.Error(err.Error())
		return ErrFindData
	}

	if data.Password != hash.Hash([]byte(currentPassword)) {
		return ErrOldPassword
	}

	{
		data.Username = ""
		data.Password = hash.Hash([]byte(newPassword))
	}

	err = service.db.Update(ctx, data)
	if err != nil {
		service.config.Logger.Error(err.Error())
		return ErrUpdateData
	}

	return nil
}

func (service *service_impl) RefreshToken(oldToken string) (token string, err error) {
	defer service.config.Logger.Info("refresh token", zap.String("execution time", executionTime(time.Now())))

	jwtData, valid := service.jwt.Parse(oldToken)
	if !valid {
		return "", ErrTokenInvalid
	}

	token, err = service.jwt.Generate(jwtData)
	if err != nil {
		service.config.Logger.Error(err.Error())
		return "", ErrGenerateToken
	}

	return token, nil
}

func (service *service_impl) UpdateProfile(ctx context.Context, token string, profile Profile) (err error) {
	defer service.config.Logger.Info("update user profile", zap.String("execution time", executionTime(time.Now())))

	jwtData, valid := service.jwt.Parse(token)
	if !valid {
		return ErrTokenInvalid
	}

	var data database.Data
	{
		data.ID, err = primitive.ObjectIDFromHex(jwtData.ID)
		if err != nil {
			service.config.Logger.Error(err.Error())
			return ErrIDInvalid
		}
	}

	filter := database.DataFilter{
		ID: true,
	}

	err = service.db.Find(ctx, &data, filter)
	if err != nil {
		service.config.Logger.Error(err.Error())
		return ErrFindData
	}

	err = service.db.Update(ctx, database.Data{
		ID:        data.ID,
		Firstname: profile.Firstname,
		Lastname:  profile.Lastname,
		Email:     profile.Email,
	})

	if err != nil {
		service.config.Logger.Error(err.Error())
		return ErrUpdateData
	}

	return nil

}

func (service *service_impl) GetProfile(ctx context.Context, token string, filter ProfileFilter) (profile Profile, err error) {
	defer service.config.Logger.Info("get user profile", zap.String("execution time", executionTime(time.Now())))

	jwtData, valid := service.jwt.Parse(token)
	if !valid {
		return Profile{}, ErrTokenInvalid
	}

	var data database.Data
	{
		data.ID, err = primitive.ObjectIDFromHex(jwtData.ID)
		if err != nil {
			service.config.Logger.Error(err.Error())
			return Profile{}, ErrIDInvalid
		}
	}

	dataFilter := database.DataFilter{
		Firstname: filter.Firstname,
		Lastname:  filter.Lastname,
		Email:     filter.Email,
	}

	err = service.db.Find(ctx, &data, dataFilter)
	if err != nil {
		service.config.Logger.Error(err.Error())
		return Profile{}, ErrFindData
	}

	if err != nil {
		return Profile{}, err
	}

	{
		profile.Firstname = data.Firstname
		profile.Lastname = data.Lastname
		profile.Email = data.Email
	}

	return profile, nil
}
