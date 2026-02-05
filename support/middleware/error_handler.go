package middleware

import (
	"net/http"

	"myapp/support/base"
	"myapp/support/logger"

	"github.com/gin-gonic/gin"
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

			// Log the error based on status code
			if statusCode >= 500 {
				logger.Error("[%s %s] %d - %s: %s", c.Request.Method, c.Request.URL.Path, statusCode, message, errMsg)
			} else if statusCode >= 400 {
				logger.Warn("[%s %s] %d - %s: %s", c.Request.Method, c.Request.URL.Path, statusCode, message, errMsg)
			}

			c.AbortWithStatusJSON(statusCode, base.CreateFailResponse(
				message,
				errMsg,
				uint(statusCode),
			))
		}
	}
}
