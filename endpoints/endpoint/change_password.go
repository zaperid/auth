package endpoint

import (
	"context"
	"montrek-auth/service"

	"github.com/go-kit/kit/endpoint"
)

type ChangePasswordRequest struct {
	Token              string `json:"token" form:"token"`
	CaptchaToken       string `json:"captcha_token" form:"captcha_token"`
	CaptchaAnswer      string `json:"captcha_answer" form:"captcha_answer"`
	CurrentPassword    string `json:"current_password" form:"current_password"`
	NewPassword        string `json:"new_password" form:"new_password"`
	NewPasswordConfirm string `json:"new_password_confirm" form:"new_password_confrim"`
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
		err = svc.ChangePassword(ctx, req.Token, req.CaptchaToken, req.CaptchaAnswer, req.CurrentPassword, req.NewPassword, req.NewPasswordConfirm)
		if err != nil {
			res.Error = err.Error()
		}

		return res, nil
	}
}
