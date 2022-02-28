package endpoints

import (
	"montrek-auth/endpoints/endpoint"
	"montrek-auth/service"
)

func NewEndpoints(svc service.Service) Endpoints {
	return Endpoints{
		GenerateCaptcha: endpoint.GenerateCaptchaEndpoint(svc),
		UsedUsername:    endpoint.UsedUsernameEndpoint(svc),
		Register:        endpoint.RegisterEndpoint(svc),
		Login:           endpoint.LoginEndpoint(svc),
	}
}
