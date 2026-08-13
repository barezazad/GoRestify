package middleware

import (
	"GoRestify/pkg/dictionary"
	"GoRestify/pkg/pkg_types"

	"github.com/gin-gonic/gin"
)

// SetDomainInHeader is used for set this domain is header
func SetDomainInHeader(app pkg_types.Enum) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("X-DOMAIN", app)
		c.Next()
	}
}

// GetLanguage is used for fetch language and ip from header
func GetLanguage() gin.HandlerFunc {
	return func(c *gin.Context) {

		lang := c.Request.Header.Get("X-LANGUAGE")
		c.Set("LANGUAGE", lang)
		if lang == "" {
			c.Set("LANGUAGE", dictionary.En)
		}

		c.Next()
	}
}
