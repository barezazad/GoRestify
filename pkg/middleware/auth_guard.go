package middleware

import (
	"strings"
	"time"

	"GoRestify/pkg/dictionary"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_jwt"
	"GoRestify/pkg/pkg_types"
	"GoRestify/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JwtAuthGuard decodes and validates an RS256 access token.
func JwtAuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			err := pkg_err.New(pkg_err.TokenIsRequired, "E7175511").
				Custom(pkg_err.UnauthorizedErr).
				Message(pkg_err.SomethingWentWrong).Build()
			response.New(c).Error(err).Abort().JSON()
			return
		}
		token := strings.TrimSpace(parts[1])

		claims, err := pkg_jwt.ParseAndValidate(token)
		if err != nil {
			checkErr(c, claims, err)
			return
		}

		if claims.TokenType != "" && claims.TokenType != pkg_jwt.TokenTypeAccess {
			err = pkg_err.New(pkg_err.TokenIsNotValid, "E7175512").
				Custom(pkg_err.UnauthorizedErr).
				Message(pkg_err.TokenIsNotValid).Build()
			response.New(c).Error(err).Abort().JSON()
			return
		}

		c.Set("USER_ID", claims.UserID)
		c.Set("USERNAME", claims.Username)
		c.Set("TOKEN", token)

		lang := c.Request.Header.Get("X-LANGUAGE")
		c.Set("LANGUAGE", lang)
		if lang == "" {
			c.Set("LANGUAGE", dictionary.En)
		}

		c.Next()
	}
}

func checkErr(c *gin.Context, claims *pkg_types.JWTClaims, err error) {
	if err == nil {
		return
	}

	switch {
	case err == jwt.ErrSignatureInvalid || errorsIsSignature(err):
		err = pkg_err.Take(err, "E1633649").Custom(pkg_err.UnauthorizedErr).
			Message(pkg_err.TokenIsNotValid).Build()
	case claims != nil && claims.ExpiresAt != nil && time.Until(claims.ExpiresAt.Time) < 10*time.Second:
		err = pkg_err.Take(err, "E1690538").Custom(pkg_err.UnauthorizedErr).
			Message(pkg_err.TokenIsExpired).Build()
	case claims != nil && claims.ExpiresAt == nil:
		err = pkg_err.Take(err, "E1652166").Custom(pkg_err.UnauthorizedErr).
			Message(pkg_err.TokenIsExpired).Build()
	default:
		err = pkg_err.Take(err, "E1655024").Custom(pkg_err.UnauthorizedErr).
			Message(pkg_err.TokenIsNotValid).Build()
	}

	response.New(c).Error(err).Abort().JSON()
}

func errorsIsSignature(err error) bool {
	return err != nil && (err.Error() == jwt.ErrSignatureInvalid.Error() ||
		strings.Contains(err.Error(), "signature is invalid") ||
		strings.Contains(err.Error(), "unexpected signing method"))
}
