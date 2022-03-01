package endpoint

import (
	"context"
	"montrek-auth/service"

	"github.com/go-kit/kit/endpoint"
)

type LoginRequest struct {
	CaptchaToken  string
	CaptchaAnswer string
	Username      string
	Password      string
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func LoginEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var res LoginResponse

		req, ok := request.(LoginRequest)
		if !ok {
			return res, ErrInvalidRequest
		}

		var err error
		res.Token, err = svc.Login(ctx, req.CaptchaToken, req.CaptchaAnswer, req.Username, req.Password)
		if err != nil {
			res.Error = err.Error()
		}

		return res, nil
	}
}
