package check_access

import (
	"GoRestify/domain/service"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_types"
	"GoRestify/pkg/response"

	"github.com/gin-gonic/gin"
)

// IsAllow will analyze if the operator should have access to special resource or not
func IsAllow(resource pkg_types.Resource) gin.HandlerFunc {
	return func(c *gin.Context) {

		var userID uint
		if id, ok := c.Get("USER_ID"); ok {
			userID = id.(uint)
		}

		accessResult := service.BaseAuthService.CheckAccess(userID, resource)

		if !accessResult {
			err := pkg_err.New("you don't have permission", "E7152019").
				Message(pkg_err.YouDontHavePermissionToThisX, resource).
				Custom(pkg_err.ForbiddenErr).Build()

			resp := response.New(c)
			resp.Error(err).Abort().JSON()
			return
		}

		c.Next()

	}
}
