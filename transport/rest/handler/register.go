package handler

import (
	"github.com/zaperid/auth/endpoints/endpoint"
	"net/http"

	"github.com/gin-gonic/gin"
	gokitEndpoint "github.com/go-kit/kit/endpoint"
)

func RegisterHandler(endpointfn gokitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req endpoint.RegisterRequest

		if c.Request.Method != http.MethodGet && c.BindJSON(&req) != nil {
			return
		} else if c.BindQuery(&req) != nil {
			return
		}

		var res endpoint.RegisterResponse
		{
			rawRes, err := endpointfn(c.Request.Context(), req)
			if err != nil {
				switch err {
				case endpoint.ErrInvalidRequest:
					c.Status(http.StatusServiceUnavailable)
					return
				default:
					c.Status(http.StatusServiceUnavailable)
					return
				}
			}

			var ok bool
			res, ok = rawRes.(endpoint.RegisterResponse)
			if !ok {
				c.Status(http.StatusBadGateway)
				return
			}
		}

		c.JSON(http.StatusOK, res)
	}
}
