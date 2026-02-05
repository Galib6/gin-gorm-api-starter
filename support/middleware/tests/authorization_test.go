package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"myapp/core/service"
	"myapp/support/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// --- Test Helpers ---

func setupAuthorizationTest(t *testing.T) (*gin.Engine, service.JWTService) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler())

	jwtS := service.NewJWTService()
	r.GET("/protected", middleware.Authenticate(jwtS), middleware.Authorize(), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	return r, jwtS
}

// --- Tests ---

func TestAuthorize_ForbiddenRole(t *testing.T) {
	r, jwtS := setupAuthorizationTest(t)

	token := jwtS.GenerateToken("abc", "user")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestAuthorize_ValidRole(t *testing.T) {
	r, jwtS := setupAuthorizationTest(t)

	token := jwtS.GenerateToken("abc", "admin")

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "ok", w.Body.String())
}
