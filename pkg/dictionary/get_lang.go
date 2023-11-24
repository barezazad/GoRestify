package dictionary

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// GetLang return suitable language according to 1.query, 2.JWT, 3.environment
func GetLang(c *gin.Context) Lang {
	var langLevel Lang

	// priority 3: get from code
	langLevel = En

	// priority 2: get from context
	langJWT, ok := c.Get("LANGUAGE")
	if ok {
		langLevel = Lang(fmt.Sprintf("%v", langJWT))
	}

	// priority 1: get from query request
	langQuery := c.Query("lang")
	if langQuery != "" {
		langLevel = Lang(langQuery)
	}

	switch langLevel {
	case En:
		return En
	case Ku:
		return Ku
	case Ar:
		return Ar
	}

	return En
}
