package endpoints

import (
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GenerateCaptcha endpoint.Endpoint
	Register        endpoint.Endpoint
	UsedUsername    endpoint.Endpoint
	Login           endpoint.Endpoint
	ChangePassword  endpoint.Endpoint
	RefreshToken    endpoint.Endpoint
	UpdateProfile   endpoint.Endpoint
	GetProfile      endpoint.Endpoint
}
