package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/zetsux/gin-gorm-clean-starter/support/base"
	"github.com/zetsux/gin-gorm-clean-starter/support/middleware"
)

// --- Test Helpers ---

func setupErrorHandlerTest(t *testing.T) *gin.Engine {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.ErrorHandler())
	return router
}

// --- Tests ---

func TestErrorHandler_NoError(t *testing.T) {
	router := setupErrorHandlerTest(t)
	router.GET("/ok", func(c *gin.Context) {
		c.JSON(http.StatusOK, base.CreateSuccessResponse("OK", http.StatusOK, nil))
	})

	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	require.Contains(t, resp.Body.String(), `"success":true`)
}

func TestErrorHandler_AppError(t *testing.T) {
	router := setupErrorHandlerTest(t)
	router.GET("/apperror", func(c *gin.Context) {
		_ = c.Error(base.NewAppError(http.StatusBadRequest, "bad request", errors.New("invalid input")))
	})

	req := httptest.NewRequest(http.MethodGet, "/apperror", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Contains(t, resp.Body.String(), `"success":false`)
	require.Contains(t, resp.Body.String(), `"bad request"`)
	require.Contains(t, resp.Body.String(), `"invalid input"`)
}

func TestErrorHandler_GenericError(t *testing.T) {
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.GET("/generic", func(c *gin.Context) {
		_ = c.Error(errors.New("something went wrong"))
	})

	req := httptest.NewRequest(http.MethodGet, "/generic", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusInternalServerError, resp.Code)
	require.Contains(t, resp.Body.String(), `"Internal Server Error"`)
	require.Contains(t, resp.Body.String(), `"something went wrong"`)
}
