package endpoint

import (
	"context"
	"montrek-auth/service"

	"github.com/go-kit/kit/endpoint"
)

type RefreshTokenRequest struct {
	Token string `json:"token" form:"token"`
}

type RefreshTokenResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func RefreshTokenEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var res LoginResponse

		req, ok := request.(RefreshTokenRequest)
		if !ok {
			return res, ErrInvalidRequest
		}

		{
			var err error
			res.Token, err = svc.RefreshToken(req.Token)
			if err != nil {
				res.Error = err.Error()
			}
		}

		return res, nil
	}
}
