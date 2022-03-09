package endpoint

import (
	"context"
	"montrek-auth/service"

	"github.com/go-kit/kit/endpoint"
)

type UpdateProfileRequest struct {
	Token     string `json:"token" form:"token"`
	Firstname string `json:"firstname" form:"firstname"`
	Lastname  string `json:"lastname" form:"lastname"`
	Email     string `json:"email" form:"email"`
}

type UpdateProfileResponse struct {
	Error string `json:"error,omitempty"`
}

func UpdateProfileEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var res UpdateProfileResponse

		req, ok := request.(UpdateProfileRequest)
		if !ok {
			return res, ErrInvalidRequest
		}

		var err error
		err = svc.UpdateProfile(ctx, req.Token, req.Firstname, req.Lastname, req.Email)
		if err != nil {
			res.Error = err.Error()
		}

		return res, nil
	}
}
