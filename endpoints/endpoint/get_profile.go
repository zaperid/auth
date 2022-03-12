package endpoint

import (
	"context"
	"github.com/zaperid/auth/service"

	"github.com/go-kit/kit/endpoint"
)

type GetProfileRequest struct {
	Token     string `json:"token" form:"token"`
	Firstname bool   `json:"firstname,omitempty" form:"firstname,omitempty"`
	Lastname  bool   `json:"lastname,omitempty" form:"lastname,omitempty"`
	Email     bool   `json:"email,omitempty" form:"email,omitempty"`
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

		filter := service.ProfileFilter{
			Firstname: req.Firstname,
			Lastname:  req.Lastname,
			Email:     req.Email,
		}

		var err error
		var profile service.Profile
		profile, err = svc.GetProfile(ctx, req.Token, filter)
		if err != nil {
			res.Error = err.Error()
		}

		{
			res.Firstname = profile.Firstname
			res.Lastname = profile.Lastname
			res.Email = profile.Email
		}

		return res, nil
	}
}
