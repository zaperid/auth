package endpoint

import (
	"context"
	"montrek-auth/service"

	"github.com/go-kit/kit/endpoint"
)

type LoginRequest struct {
	CaptchaToken  string `json:"captcha_token" form:"captcha_token"`
	CaptchaAnswer string `json:"captcha_answer" form:"captcha_answer"`
	Username      string `json:"username" form:"captcha_username"`
	Password      string `json:"password" form:"capthca_answer"`
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
