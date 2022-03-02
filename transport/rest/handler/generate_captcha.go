package handler

import (
	"montrek-auth/endpoints/endpoint"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	gokitEndpoint "github.com/go-kit/kit/endpoint"
)

func GenerateCaptchaHandler(endpointfn gokitEndpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req endpoint.GenerateCaptchaRequest

		{
			if c.Request.Method != http.MethodGet && c.BindJSON(&req) != nil {
				return
			} else {
				{
					var err error

					heighStr := c.Query("height")
					req.Height, err = strconv.Atoi(heighStr)
					if err != nil {
						c.Status(http.StatusBadRequest)
						return
					}
				}

				{
					var err error

					widthStr := c.Query("width")
					req.Width, err = strconv.Atoi(widthStr)
					if err != nil {
						c.Status(http.StatusBadRequest)
						return
					}
				}
			}
		}

		var res endpoint.GenerateCaptchaResponse
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
			res, ok = rawRes.(endpoint.GenerateCaptchaResponse)
			if !ok {
				c.Status(http.StatusBadGateway)
				return
			}
		}

		c.JSON(http.StatusOK, res)
	}
}
