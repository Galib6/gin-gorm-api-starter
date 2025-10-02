package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zetsux/gin-gorm-clean-starter/common/base"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/dto"
	"github.com/zetsux/gin-gorm-clean-starter/tests/support"
)

// Test update user name endpoint
func TestIntegration_UpdateUserName(t *testing.T) {
	testApp := support.SetupTestApp(t)
	server := testApp.Server

	// Create user and get token
	token := createUserAndGetToken(t, server, "Bob Wilson", "bob@example.com", "password123")

	// Test update name
	updateReq := dto.UserNameUpdateRequest{Name: "Bob Updated"}
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/me/name", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp base.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.True(t, resp.IsSuccess)
	require.Equal(t, "User update successful", resp.Message)

	// Verify the response
	userData := resp.Data.(map[string]interface{})
	require.Equal(t, "Bob Updated", userData["name"])

	// Fetch profile to verify
	reqMe := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	reqMe.Header.Set("Authorization", "Bearer "+token)
	wMe := httptest.NewRecorder()

	server.ServeHTTP(wMe, reqMe)

	require.Equal(t, http.StatusOK, wMe.Code)

	var respMe base.Response
	err = json.Unmarshal(wMe.Body.Bytes(), &respMe)
	require.NoError(t, err)
	require.True(t, respMe.IsSuccess)

	// Verify the updated data
	userData = respMe.Data.(map[string]interface{})
	require.Equal(t, "Bob Updated", userData["name"])
	require.Equal(t, "bob@example.com", userData["email"])
}

// Test update user name without authentication
func TestIntegration_UpdateUserName_Unauthorized(t *testing.T) {
	testApp := support.SetupTestApp(t)
	server := testApp.Server

	updateReq := dto.UserNameUpdateRequest{Name: "New Name"}
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/me/name", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)

	var resp base.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.False(t, resp.IsSuccess)
}

// Test multiple users isolation
func TestIntegration_MultipleUsersIsolation(t *testing.T) {
	testApp := support.SetupTestApp(t)
	server := testApp.Server

	// Create multiple users
	token1 := createUserAndGetToken(t, server, "User One", "user1@example.com", "pass1")
	token2 := createUserAndGetToken(t, server, "User Two", "user2@example.com", "pass2")

	// User 1 updates their name
	updateReq := dto.UserNameUpdateRequest{Name: "User One Updated"}
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/me/name", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token1)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// Verify User 2 is unaffected
	req = httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token2)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var resp base.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	userData := resp.Data.(map[string]interface{})
	require.Equal(t, "User Two", userData["name"]) // Should be unchanged
	require.Equal(t, "user2@example.com", userData["email"])

	// Verify User 1 has the updated name
	req = httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token1)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	json.Unmarshal(w.Body.Bytes(), &resp)
	userData = resp.Data.(map[string]interface{})
	require.Equal(t, "User One Updated", userData["name"])
	require.Equal(t, "user1@example.com", userData["email"])
}
