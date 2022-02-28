package endpoint

import (
	"context"
	"montrek-auth/service"

	"github.com/go-kit/kit/endpoint"
)

type UsedUsernameRequest struct {
	Username string
}

type UsedUsernameResponse struct {
	Used  bool   `json:"used,omitempty"`
	Error string `json:"error,omitempty"`
}

func UsedUsernameEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var res UsedUsernameResponse

		req, ok := request.(UsedUsernameRequest)
		if !ok {
			return res, ErrInvalidRequest
		}

		var err error
		res.Used, err = svc.UsedUsername(ctx, req.Username)
		if err != nil {
			res.Error = err.Error()
		}

		return res, nil
	}
}
