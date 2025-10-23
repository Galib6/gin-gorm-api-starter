package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zetsux/gin-gorm-api-starter/core/helper/dto"
	"github.com/zetsux/gin-gorm-api-starter/core/helper/messages"
	"github.com/zetsux/gin-gorm-api-starter/support/base"
	"github.com/zetsux/gin-gorm-api-starter/support/constant"
	"github.com/zetsux/gin-gorm-api-starter/tests/testutil"
	"github.com/zetsux/gin-gorm-api-starter/tests/testutil/factory"
)

// Test get users by admin endpoint
func TestIntegration_GetUsersByAdmin(t *testing.T) {
	testApp := testutil.SetupTestApp(t)
	server := testApp.Server

	// Create admin user and get token
	adminEmail, adminPass := "admin@mail.com", "password123"
	factory.SeedUser(t, testApp.UserRepo, "Admin User", adminEmail, adminPass, constant.EnumRoleAdmin)
	token := testutil.GetToken(t, server, adminEmail, adminPass)

	// Create regular users
	factory.SeedUsers(t, testApp.UserRepo, 15)

	// Test get users
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp base.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.True(t, resp.IsSuccess)
	require.Equal(t, messages.MsgUsersFetchSuccess, resp.Message)

	// Verify response data contains user info
	userData := resp.Data.([]interface{})
	require.NotEmpty(t, userData)
	require.Equal(t, len(userData), 16) // 15 regular + 1 admin
}

// Test get users by non-admin
func TestIntegration_GetUsersByAdmin_Forbidden(t *testing.T) {
	testApp := testutil.SetupTestApp(t)
	server := testApp.Server

	// Create regular user and get token
	token := testutil.CreateUserAndGetToken(t, server, "Regular User", "regular@example.com", "password123")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)

	var resp base.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.False(t, resp.IsSuccess)
}

// Test delete user by admin endpoint
func TestIntegration_DeleteUserByAdmin(t *testing.T) {
	testApp := testutil.SetupTestApp(t)
	server := testApp.Server

	// Create admin user and get token
	adminEmail, adminPass := "admin@mail.com", "password123"
	factory.SeedUser(t, testApp.UserRepo, "Admin User", adminEmail, adminPass, constant.EnumRoleAdmin)
	token := testutil.GetToken(t, server, adminEmail, adminPass)

	// Create regular user to delete
	user := factory.SeedUser(t, testApp.UserRepo, "User To Delete", "user_to_delete@example.com", "password123", constant.EnumRoleUser)

	// Test delete user
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+user.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp base.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.True(t, resp.IsSuccess)
	require.Equal(t, messages.MsgUserDeleteSuccess, resp.Message)

	// Verify user cannot login after deletion
	loginReq := dto.UserLoginRequest{
		Email:    user.Email,
		Password: user.Password,
	}

	body, _ := json.Marshal(loginReq)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	server.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.False(t, resp.IsSuccess)
}

// Test delete user by non-admin
func TestIntegration_DeleteUserByAdmin_Forbidden(t *testing.T) {
	testApp := testutil.SetupTestApp(t)
	server := testApp.Server

	// Create regular user and get token
	token := testutil.CreateUserAndGetToken(t, server, "Regular User", "regular@example.com", "password123")

	// Create another user to attempt deletion
	user := factory.SeedUser(t, testApp.UserRepo, "User To Delete", "user_to_delete@example.com", "password123", constant.EnumRoleUser)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+user.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)

	var resp base.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.False(t, resp.IsSuccess)
}

// Test update user by admin endpoint
func TestIntegration_UpdateUserByAdmin(t *testing.T) {
	testApp := testutil.SetupTestApp(t)
	server := testApp.Server

	// Create admin user and get token
	adminEmail, adminPass := "admin@mail.com", "password123"
	factory.SeedUser(t, testApp.UserRepo, "Admin User", adminEmail, adminPass, constant.EnumRoleAdmin)
	token := testutil.GetToken(t, server, adminEmail, adminPass)

	// Create regular user to update
	user := factory.SeedUser(t, testApp.UserRepo, "User To Update", "user_to_update@example.com", "password123", constant.EnumRoleUser)

	// Test update user
	reqBody := dto.UserUpdateRequest{
		Name:  "Updated User",
		Email: "updated_user@example.com",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/"+user.ID.String(), bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp base.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.True(t, resp.IsSuccess)
	require.Equal(t, messages.MsgUserUpdateSuccess, resp.Message)

	// Verify response data contains updated info
	userData := resp.Data.(map[string]interface{})
	require.Equal(t, "Updated User", userData["name"])
	require.Equal(t, "updated_user@example.com", userData["email"])

	// Verify user can login with updated email
	_ = testutil.GetToken(t, server, "updated_user@example.com", "password123")
}

// Test update user by non-admin
func TestIntegration_UpdateUserByAdmin_Forbidden(t *testing.T) {
	testApp := testutil.SetupTestApp(t)
	server := testApp.Server

	// Create regular user and get token
	token := testutil.CreateUserAndGetToken(t, server, "Regular User", "regular@example.com", "password123")

	// Create another user to attempt update
	user := factory.SeedUser(t, testApp.UserRepo, "User To Update", "user_to_update@example.com", "password123", constant.EnumRoleUser)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/"+user.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)

	var resp base.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.False(t, resp.IsSuccess)
}
