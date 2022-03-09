package endpoint

import (
	"context"
	"montrek-auth/service"

	"github.com/go-kit/kit/endpoint"
)

type GetProfileRequest struct {
	Token string `json:"token" form:"token"`
}

type GetProfileResponse struct {
	Error     string `json:"error,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Email     string `json:"email,omitempty"`
}

func GetProfileEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var res GetProfileResponse

		req, ok := request.(GetProfileRequest)
		if !ok {
			return res, ErrInvalidRequest
		}

		var err error
		res.Firstname, res.Lastname, res.Email, err = svc.GetProfile(ctx, req.Token)
		if err != nil {
			res.Error = err.Error()
		}

		return res, nil
	}
}
