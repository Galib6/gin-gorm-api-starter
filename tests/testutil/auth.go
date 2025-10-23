package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/zetsux/gin-gorm-api-starter/core/helper/dto"
	"github.com/zetsux/gin-gorm-api-starter/support/base"
)

func GetToken(t *testing.T, server *gin.Engine, email, password string) string {
	t.Helper()

	// Login to get token
	loginReq := dto.UserLoginRequest{
		Email:    email,
		Password: password,
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var loginResp base.Response
	err := json.Unmarshal(w.Body.Bytes(), &loginResp)
	require.NoError(t, err)

	authData := loginResp.Data.(map[string]interface{})
	token := authData["token"].(string)
	require.NotEmpty(t, token)

	return token
}

func CreateUserAndGetToken(t *testing.T, server *gin.Engine, name, email, password string) string {
	t.Helper()

	// Register user
	regReq := dto.UserRegisterRequest{
		Name:     name,
		Email:    email,
		Password: password,
	}

	body, _ := json.Marshal(regReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	return GetToken(t, server, email, password)
}
