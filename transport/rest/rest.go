package rest

import (
	"fmt"
	"montrek-auth/endpoints"
	"montrek-auth/transport/rest/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type server_impl struct {
	endpoints endpoints.Endpoints
	config    Config
	gin       *gin.Engine
}

func NewServer(endpoints endpoints.Endpoints, config Config) Server {
	gin.SetMode(gin.ReleaseMode)

	server := server_impl{
		endpoints: endpoints,
		config:    config,
		gin:       gin.New(),
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.AllowOrigins
	server.gin.Use(cors.New(corsConfig))

	server.setup()

	return &server
}

func (server *server_impl) Run() error {
	host := fmt.Sprintf(":%d", server.config.Port)
	server.config.Logger.Info("start server", zap.String("host", host))

	return server.gin.Run(host)
}

func (server *server_impl) setup() {
	server.gin.GET("/api/v1/captcha", handler.GenerateCaptchaHandler(server.endpoints.GenerateCaptcha))

	server.gin.POST("/api/v1/users", handler.RegisterHandler(server.endpoints.Register))
	server.gin.PUT("/api/v1/users/password", handler.ChangePasswordHandler(server.endpoints.ChangePassword))
	server.gin.GET("/api/v1/users/profile", handler.GetProfile(server.endpoints.GetProfile))
	server.gin.PUT("/api/v1/users/profile", handler.UpdateProfileHandler(server.endpoints.UpdateProfile))
	server.gin.GET("/api/v1/users/token", handler.RefreshTokenHandler(server.endpoints.RefreshToken))
	server.gin.POST("/api/v1/users/token", handler.LoginHandler(server.endpoints.Login))
	server.gin.GET("/api/v1/users/username", handler.UsedUsernameHandler(server.endpoints.UsedUsername))
}
