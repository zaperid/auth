package endpoint

import (
	"context"
	"montrek-auth/service"

	"github.com/go-kit/kit/endpoint"
)

type ChangePasswordRequest struct {
	Token              string
	CaptchaToken       string
	CaptchaAnswer      string
	OldPassword        string
	NewPassword        string
	NewPasswordConfirm string
}

type ChangePasswordResponse struct {
	Error string `json:"error,omitempty"`
}

func ChangePasswordEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var res ChangePasswordResponse

		req, ok := request.(ChangePasswordRequest)
		if !ok {
			return res, ErrInvalidRequest
		}

		var err error
		err = svc.ChangePassword(ctx, req.Token, req.CaptchaToken, req.CaptchaAnswer, req.OldPassword, req.NewPassword, req.NewPasswordConfirm)
		if err != nil {
			res.Error = err.Error()
		}

		return res, nil
	}
}
