package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"myapp/support/base"
	"myapp/support/logger"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware recovers from panics and returns a 500 error
// This is Go's equivalent of exception filtering/global exception handler
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get stack trace
				stack := debug.Stack()
				stackLines := strings.Split(string(stack), "\n")

				// Log the panic with stack trace
				logger.Error("🔥 PANIC RECOVERED: %v", err)
				logger.Error("Stack trace:")
				for _, line := range stackLines {
					if line != "" {
						fmt.Printf("    %s\n", line)
					}
				}

				// Return error response
				c.AbortWithStatusJSON(http.StatusInternalServerError, base.CreateFailResponse(
					"Internal Server Error",
					"An unexpected error occurred",
					http.StatusInternalServerError,
				))
			}
		}()
		c.Next()
	}
}

// RequestLoggerMiddleware logs HTTP requests with colored output
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get request details
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		// Log the request
		logger.LogHTTPRequest(method, path, statusCode, latency, clientIP)
	}
}
