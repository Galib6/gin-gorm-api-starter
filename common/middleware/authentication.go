package middleware

import (
	"net/http"
	"strings"

	"github.com/zetsux/gin-gorm-clean-starter/common/base"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/messages"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"

	"github.com/gin-gonic/gin"
)

func Authenticate(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			_ = c.Error(base.NewAppError(http.StatusUnauthorized,
				messages.MsgAuthNoToken, nil))
			c.Abort()
			return
		}
		if !strings.Contains(authHeader, "Bearer ") {
			_ = c.Error(base.NewAppError(http.StatusUnauthorized,
				messages.MsgAuthNoToken, nil))
			c.Abort()
			return
		}
		authHeader = strings.ReplaceAll(authHeader, "Bearer ", "")
		token, err := jwtService.ValidateToken(authHeader)
		if err != nil {
			_ = c.Error(base.NewAppError(http.StatusUnauthorized,
				messages.MsgAuthInvalidToken, err))
			c.Abort()
			return
		}

		if !token.Valid {
			_ = c.Error(base.NewAppError(http.StatusUnauthorized,
				messages.MsgAuthInvalidToken, nil))
			c.Abort()
			return
		}

		// get role from token
		idRes, roleRes, err := jwtService.GetAttrByToken(authHeader)
		if err != nil {
			_ = c.Error(base.NewAppError(http.StatusUnauthorized,
				messages.MsgAuthFailedProcess, err))
			c.Abort()
			return
		}
		c.Set("ID", idRes)
		c.Set("ROLE", roleRes)
		c.Next()
	}
}
