package handler

import (
	"github.com/zaperid/auth/endpoints/endpoint"
	"net/http"

	"github.com/gin-gonic/gin"
	gokitendpoint "github.com/go-kit/kit/endpoint"
)

func RefreshTokenHandler(endpointfn gokitendpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req endpoint.RefreshTokenRequest

		if c.Request.Method != http.MethodGet && c.BindJSON(&req) != nil {
			return
		} else if c.BindQuery(&req) != nil {
			return
		}

		res, err := endpointfn(c.Request.Context(), req)
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

		c.JSON(http.StatusOK, res)
	}
}
