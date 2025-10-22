package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zetsux/gin-gorm-api-starter/support/base"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			statusCode := http.StatusInternalServerError
			message := "Internal Server Error"
			errMsg := ""

			if appErr, ok := err.(*base.AppError); ok {
				statusCode = appErr.Code
				message = appErr.Message

				if appErr.Err != nil {
					errMsg = appErr.Err.Error()
				}
			} else {
				errMsg = err.Error()
			}

			c.AbortWithStatusJSON(statusCode, base.CreateFailResponse(
				message,
				errMsg,
				uint(statusCode),
			))
		}
	}
}
