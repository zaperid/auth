package handler

import (
	"fmt"
	"github.com/zaperid/auth/endpoints/endpoint"
	"net/http"

	"github.com/gin-gonic/gin"
	gokitEndpoint "github.com/go-kit/kit/endpoint"
)

func GetProfile(endpointfn gokitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req endpoint.GetProfileRequest

		if c.Request.Method != http.MethodGet && c.BindJSON(&req) != nil {
			return
		} else if c.BindQuery(&req) != nil {
			return
		}

		fmt.Println(req)

		var res endpoint.GetProfileResponse
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
			res, ok = rawRes.(endpoint.GetProfileResponse)
			if !ok {
				c.Status(http.StatusBadGateway)
				return
			}

		}

		c.JSON(http.StatusOK, res)
	}
}
