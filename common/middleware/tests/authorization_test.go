package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/zetsux/gin-gorm-clean-starter/common/middleware"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"
)

func TestAuthenticate_ForbiddenRole(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	jwtS := service.NewJWTService()
	r.GET("/protected", middleware.Authenticate(jwtS), middleware.Authorize("admin"), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	token := jwtS.GenerateToken("abc", "user")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusForbidden, w.Code)
}
