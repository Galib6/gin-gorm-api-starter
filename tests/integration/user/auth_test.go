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
	"github.com/zetsux/gin-gorm-api-starter/tests/testutil"
)

// Test user registration endpoint
func TestIntegration_UserRegistration(t *testing.T) {
	testApp := testutil.SetupTestApp(t)
	server := testApp.Server

	regReq := dto.UserRegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(regReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var regResp base.Response
	err := json.Unmarshal(w.Body.Bytes(), &regResp)
	require.NoError(t, err)
	require.True(t, regResp.IsSuccess)
	require.Equal(t, messages.MsgUserRegisterSuccess, regResp.Message)

	// Verify response data contains user info
	userData := regResp.Data.(map[string]interface{})
	require.Equal(t, regReq.Name, userData["name"])
	require.Equal(t, regReq.Email, userData["email"])
	require.Equal(t, "user", userData["role"])
	require.NotEmpty(t, userData["id"])
}

// Test user registration with invalid data
func TestIntegration_UserRegistration_InvalidData(t *testing.T) {
	testApp := testutil.SetupTestApp(t)
	server := testApp.Server

	testCases := []struct {
		name        string
		request     dto.UserRegisterRequest
		expectedMsg string
	}{
		{
			name: "missing email",
			request: dto.UserRegisterRequest{
				Name:     "John Doe",
				Password: "password123",
			},
		},
		{
			name: "missing password",
			request: dto.UserRegisterRequest{
				Name:  "John Doe",
				Email: "john@example.com",
			},
		},
		{
			name: "missing name",
			request: dto.UserRegisterRequest{
				Email:    "john@example.com",
				Password: "password123",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.request)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)

			var resp base.Response
			json.Unmarshal(w.Body.Bytes(), &resp)
			require.False(t, resp.IsSuccess)
		})
	}
}

// Test duplicate email registration
func TestIntegration_UserRegistration_DuplicateEmail(t *testing.T) {
	testApp := testutil.SetupTestApp(t)
	server := testApp.Server

	regReq := dto.UserRegisterRequest{
		Name:     "John Doe",
		Email:    "duplicate@example.com",
		Password: "password123",
	}

	// First registration should succeed
	body, _ := json.Marshal(regReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// Second registration with same email should fail
	req = httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp base.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.False(t, resp.IsSuccess)
}

// Test user login endpoint
func TestIntegration_UserLogin(t *testing.T) {
	testApp := testutil.SetupTestApp(t)
	server := testApp.Server

	// First register a user
	regReq := dto.UserRegisterRequest{
		Name:     "Jane Smith",
		Email:    "jane@example.com",
		Password: "secret123",
	}

	body, _ := json.Marshal(regReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// Now test login
	loginReq := dto.UserLoginRequest{
		Email:    regReq.Email,
		Password: regReq.Password,
	}

	body, _ = json.Marshal(loginReq)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	server.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var loginResp base.Response
	err := json.Unmarshal(w.Body.Bytes(), &loginResp)
	require.NoError(t, err)
	require.True(t, loginResp.IsSuccess)
	require.Equal(t, messages.MsgUserLoginSuccess, loginResp.Message)

	// Verify response contains token and user data
	authData := loginResp.Data.(map[string]interface{})
	token := authData["token"].(string)
	require.NotEmpty(t, token)

	role := authData["role"].(string)
	require.Equal(t, role, "user")
}

// Test user login with invalid credentials
func TestIntegration_UserLogin_InvalidCredentials(t *testing.T) {
	testApp := testutil.SetupTestApp(t)
	server := testApp.Server

	testCases := []struct {
		name     string
		email    string
		password string
	}{
		{
			name:     "nonexistent user",
			email:    "nonexistent@example.com",
			password: "password123",
		},
		{
			name:     "wrong password",
			email:    "test@example.com",
			password: "wrongpassword",
		},
	}

	// Create a test user for wrong password test
	testutil.CreateUserAndGetToken(t, server, "Test User", "test@example.com", "correctpassword")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			loginReq := dto.UserLoginRequest{
				Email:    tc.email,
				Password: tc.password,
			}

			body, _ := json.Marshal(loginReq)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)

			require.Equal(t, http.StatusUnauthorized, w.Code)

			var resp base.Response
			json.Unmarshal(w.Body.Bytes(), &resp)
			require.False(t, resp.IsSuccess)
		})
	}
}
