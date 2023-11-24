package middleware

import (
	"time"

	"GoRestify/pkg/dictionary"
	"GoRestify/pkg/pkg_config"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_types"
	"GoRestify/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JwtAuthGuard is used for decode the token and get public and private information
func JwtAuthGuard() gin.HandlerFunc {

	jwtKey := []byte(pkg_config.Config.JWTSecretKey)

	fJWT := func(token *jwt.Token) (any, error) { return jwtKey, nil }

	return func(c *gin.Context) {

		tokenArr, ok := c.Request.Header["Authorization"]
		if !ok || len(tokenArr[0]) <= 7 {
			err := pkg_err.New(pkg_err.TokenIsRequired, "E7175511").
				Custom(pkg_err.UnauthorizedErr).
				Message(pkg_err.SomethingWentWrong).Build()
			response.New(c).Error(err).Abort().JSON()
			return
		}
		token := tokenArr[0][7:]

		claims := &pkg_types.JWTClaims{}

		if tkn, err := jwt.ParseWithClaims(token, claims, fJWT); err != nil {
			checkErr(c, claims, err)
			return
		} else if !tkn.Valid {
			checkToken(c, tkn)
			return
		}

		c.Set("USER_ID", claims.UserID)
		c.Set("USERNAME", claims.Username)
		c.Set("TOKEN", token)
		c.Set("EMAIL", claims.Email)
		c.Set("PHONE", claims.Phone)

		// set language
		lang := c.Request.Header.Get("X-LANGUAGE")
		c.Set("LANGUAGE", lang)
		if lang == "" {
			c.Set("LANGUAGE", dictionary.En)
		}

		c.Next()
	}
}

func checkErr(c *gin.Context, claims *pkg_types.JWTClaims, err error) {
	if err != nil {
		switch {
		case err == jwt.ErrSignatureInvalid:
			err = pkg_err.Take(err, "E1633649").Custom(pkg_err.UnauthorizedErr).
				Message(pkg_err.TokenIsNotValid).Build()
			response.New(c).Error(err).Abort().JSON()
			return

		case claims.RegisteredClaims.ExpiresAt == nil:
			err = pkg_err.Take(err, "E1652166").Custom(pkg_err.UnauthorizedErr).
				Message(pkg_err.TokenIsExpired).Build()
			response.New(c).Error(err).Abort().JSON()
			return

		case time.Until(claims.RegisteredClaims.ExpiresAt.Time) < 10*time.Second:
			err = pkg_err.Take(err, "E1690538").Custom(pkg_err.UnauthorizedErr).
				Message(pkg_err.TokenIsExpired).Build()
			response.New(c).Error(err).Abort().JSON()
			return

		default:
			err = pkg_err.Take(err, "E1655024").Custom(pkg_err.UnauthorizedErr).
				Message(pkg_err.TokenIsNotValid).Build()
			response.New(c).Error(err).Abort().JSON()
			return
		}
	}
}

func checkToken(c *gin.Context, token *jwt.Token) {
	if !token.Valid {
		err := pkg_err.New(pkg_err.TokenIsNotValid, "E7184392").Custom(pkg_err.UnauthorizedErr).
			Message(pkg_err.TokenIsNotValid).Build()
		response.New(c).Error(err).Abort().JSON()
		return
	}
}
