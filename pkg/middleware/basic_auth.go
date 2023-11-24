package middleware

import (
	"encoding/base64"
	"strings"

	"GoRestify/pkg/pkg_config"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/response"

	"github.com/gin-gonic/gin"
)

// BasicAuthGuard is used for decode the basic auth in header for third party apps
func BasicAuthGuard() gin.HandlerFunc {

	return func(c *gin.Context) {

		var payload []byte

		basicAuth := strings.TrimSpace(c.Query("BasicAuth"))
		if basicAuth != "" {
			payload, _ = base64.StdEncoding.DecodeString(basicAuth)
		} else {

			basicAuth := strings.SplitN(c.Request.Header.Get("X-Authorization"), " ", 2)
			if len(basicAuth) != 2 || basicAuth[0] != "Basic" {
				err := pkg_err.New(pkg_err.BasicAuthIsRequired, "E1130018").
					Custom(pkg_err.UnauthorizedErr).
					Message(pkg_err.BasicAuthIsRequired).Build()
				response.New(c).Error(err).Abort().JSON()
				return
			}
			payload, _ = base64.StdEncoding.DecodeString(basicAuth[1])
		}

		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 || !(pkg_config.Config.BasicAuthUsername == pair[0] && pkg_config.Config.BasicAuthPassword == pair[1]) {
			err := pkg_err.New(pkg_err.BasicAuthInvalid, "E1165485").
				Custom(pkg_err.UnauthorizedErr).
				Message(pkg_err.BasicAuthInvalid).Build()
			response.New(c).Error(err).Abort().JSON()
			return
		}

		c.Next()
	}
}
