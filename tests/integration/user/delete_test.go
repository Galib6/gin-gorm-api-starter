package user_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/errors"
	"github.com/zetsux/gin-gorm-clean-starter/support/base"
	"github.com/zetsux/gin-gorm-clean-starter/tests/testutil"
)

// Test delete user endpoint
func TestIntegration_DeleteUser(t *testing.T) {
	testApp := testutil.SetupTestApp(t)
	server := testApp.Server

	// Create user and get token
	token := createUserAndGetToken(t, server, "Charlie Brown", "charlie@example.com", "password123")

	// Test delete user
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp base.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.True(t, resp.IsSuccess)
	require.Equal(t, "User delete successful", resp.Message)

	// Verify user cannot access profile after deletion
	req = httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)

	// The repository returns error
	require.Equal(t, http.StatusBadRequest, w.Code)

	json.Unmarshal(w.Body.Bytes(), &resp)
	require.False(t, resp.IsSuccess)
	require.Equal(t, errors.ErrUserNotFound.Error(), resp.Error)
}

// Test delete user without authentication
func TestIntegration_DeleteUser_Unauthorized(t *testing.T) {
	testApp := testutil.SetupTestApp(t)
	server := testApp.Server

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/me", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)

	var resp base.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.False(t, resp.IsSuccess)
}
