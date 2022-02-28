package endpoints

import (
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GenerateCaptcha endpoint.Endpoint
	Register        endpoint.Endpoint
	UsedUsername    endpoint.Endpoint
	Login           endpoint.Endpoint
}
