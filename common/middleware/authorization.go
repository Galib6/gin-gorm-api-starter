package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/zetsux/gin-gorm-clean-starter/common/base"
	"github.com/zetsux/gin-gorm-clean-starter/common/constant"
)

func Authorize(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleRes, exists := c.Get("ROLE")
		if !exists {
			response := base.CreateFailResponse("Failed to process request", "", http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		role, ok := roleRes.(string)
		if !ok {
			response := base.CreateFailResponse("Failed to process request", "", http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		if role != constant.EnumRoleAdmin && !slices.Contains(roles, role) {
			response := base.CreateFailResponse("Action unauthorized", "", http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		c.Next()
	}
}
