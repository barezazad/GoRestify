package middleware

import (
	"crypto/subtle"
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

		basicAuth := strings.SplitN(c.Request.Header.Get("X-Authorization"), " ", 2)
		if len(basicAuth) != 2 || !strings.EqualFold(basicAuth[0], "Basic") || strings.TrimSpace(basicAuth[1]) == "" {
			err := pkg_err.New(pkg_err.BasicAuthIsRequired, "E1130018").
				Custom(pkg_err.UnauthorizedErr).
				Message(pkg_err.BasicAuthIsRequired).Build()
			response.New(c).Error(err).Abort().JSON()
			return
		}

		payload, errDecode := base64.StdEncoding.DecodeString(basicAuth[1])
		if errDecode != nil {
			err := pkg_err.New(pkg_err.BasicAuthInvalid, "E1165484").
				Custom(pkg_err.UnauthorizedErr).
				Message(pkg_err.BasicAuthInvalid).Build()
			response.New(c).Error(err).Abort().JSON()
			return
		}

		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 || !basicAuthCredentialsMatch(pair[0], pair[1]) {
			err := pkg_err.New(pkg_err.BasicAuthInvalid, "E1165485").
				Custom(pkg_err.UnauthorizedErr).
				Message(pkg_err.BasicAuthInvalid).Build()
			response.New(c).Error(err).Abort().JSON()
			return
		}

		c.Next()
	}
}

func basicAuthCredentialsMatch(username, password string) bool {
	expectedUser := pkg_config.Config.BasicAuthUsername
	expectedPass := pkg_config.Config.BasicAuthPassword

	// Always compare both so timing does not reveal which field failed.
	userOK := subtle.ConstantTimeCompare([]byte(username), []byte(expectedUser)) == 1
	passOK := subtle.ConstantTimeCompare([]byte(password), []byte(expectedPass)) == 1
	return userOK && passOK
}
