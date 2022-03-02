package endpoint

import (
	"context"
	"montrek-auth/service"

	"github.com/go-kit/kit/endpoint"
)

type GenerateCaptchaRequest struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type GenerateCaptchaResponse struct {
	CaptchaToken string `json:"captcha_token,omitempty"`
	Image        string `json:"image,omitempty"`
	Error        string `json:"error,omitempty"`
}

func GenerateCaptchaEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var res GenerateCaptchaResponse

		req, ok := request.(GenerateCaptchaRequest)
		if !ok {
			return res, ErrInvalidRequest
		}

		var err error
		res.CaptchaToken, res.Image, err = svc.GenerateCaptcha(req.Height, req.Width)
		if err != nil {
			res.Error = err.Error()
		}

		return res, nil
	}
}
