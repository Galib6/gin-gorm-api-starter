package middleware

import (
	"net/http"
	"slices"

	"myapp/core/helper/messages"
	"myapp/support/base"
	"myapp/support/constant"

	"github.com/gin-gonic/gin"
)

func Authorize(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleRes, exists := c.Get("ROLE")
		if !exists {
			_ = c.Error(base.NewAppError(http.StatusUnauthorized,
				messages.MsgAuthFailedProcess, nil))
			c.Abort()
			return
		}

		role, ok := roleRes.(string)
		if !ok {
			_ = c.Error(base.NewAppError(http.StatusUnauthorized,
				messages.MsgAuthFailedProcess, nil))
			c.Abort()
			return
		}

		if role != constant.EnumRoleAdmin && !slices.Contains(roles, role) {
			_ = c.Error(base.NewAppError(http.StatusForbidden,
				messages.MsgAuthActionUnauthorized, nil))
			c.Abort()
			return
		}

		c.Next()
	}
}
