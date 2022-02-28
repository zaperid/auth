package endpoint

import (
	"context"
	"montrek-auth/service"

	"github.com/go-kit/kit/endpoint"
)

type GenerateCaptchaRequest struct {
	Height int
	Width  int
}

type GenerateCaptchaResponse struct {
	Token string `json:"token,omitempty"`
	Image string `json:"image,omitempty"`
	Error string `json:"error,omitempty"`
}

func GenerateCaptchaEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var res GenerateCaptchaResponse

		req, ok := request.(GenerateCaptchaRequest)
		if !ok {
			return res, ErrInvalidRequest
		}

		var err error
		res.Token, res.Image, err = svc.GenerateCaptcha(req.Height, req.Width)
		if err != nil {
			res.Error = err.Error()
		}

		return res, nil
	}
}