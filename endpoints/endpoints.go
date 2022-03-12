package endpoints

import (
	"github.com/zaperid/auth/endpoints/endpoint"
	"github.com/zaperid/auth/service"
)

func NewEndpoints(svc service.Service) Endpoints {
	return Endpoints{
		GenerateCaptcha: endpoint.GenerateCaptchaEndpoint(svc),
		UsedUsername:    endpoint.UsedUsernameEndpoint(svc),
		Register:        endpoint.RegisterEndpoint(svc),
		Login:           endpoint.LoginEndpoint(svc),
		ChangePassword:  endpoint.ChangePasswordEndpoint(svc),
		RefreshToken:    endpoint.RefreshTokenEndpoint(svc),
		UpdateProfile:   endpoint.UpdateProfileEndpoint(svc),
		GetProfile:      endpoint.GetProfileEndpoint(svc),
	}
}
