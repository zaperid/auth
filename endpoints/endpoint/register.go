package endpoint

import (
	"context"
	"montrek-auth/service"

	"github.com/go-kit/kit/endpoint"
)

type RegisterRequest struct {
	CaptchaToken    string `json:"captcha_token" form:"captcha_token"`
	CaptchaAnswer   string `json:"captcha_answer" form:"captcha_answer"`
	Username        string `json:"username" form:"username"`
	Password        string `json:"password" form:"password"`
	PasswordConfirm string `json:"password_confirm" form:"password_confirm"`
}

type RegisterResponse struct {
	Error string `json:"error,omitempty"`
}

func RegisterEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var res RegisterResponse

		req, ok := request.(RegisterRequest)
		if !ok {
			return res, ErrInvalidRequest
		}

		var err error
		err = svc.Register(ctx, req.CaptchaToken, req.CaptchaAnswer, req.Username, req.Password, req.PasswordConfirm)
		if err != nil {
			res.Error = err.Error()
		}

		return res, nil
	}
}
