package user_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zetsux/gin-gorm-clean-starter/support/base"
	"github.com/zetsux/gin-gorm-clean-starter/tests/testutil"
)

// Test get me endpoint
func TestIntegration_GetMe(t *testing.T) {
	testApp := testutil.SetupTestApp(t)
	server := testApp.Server

	// Create user and get token
	name := "Alice Johnson"
	email := "alice@example.com"
	password := "password123"
	token := createUserAndGetToken(t, server, name, email, password)

	// Test get me
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp base.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.True(t, resp.IsSuccess)
	require.Equal(t, "User fetched successfully", resp.Message)

	userData := resp.Data.(map[string]interface{})
	require.Equal(t, email, userData["email"])
	require.Equal(t, name, userData["name"])
	require.Equal(t, "user", userData["role"])
	require.NotEmpty(t, userData["id"])
}

// Test get me without authentication
func TestIntegration_GetMe_Unauthorized(t *testing.T) {
	testApp := testutil.SetupTestApp(t)
	server := testApp.Server

	testCases := []struct {
		name   string
		header string
	}{
		{
			name:   "no authorization header",
			header: "",
		},
		{
			name:   "invalid token",
			header: "Bearer invalid-token",
		},
		{
			name:   "malformed header",
			header: "invalid-format",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)

			require.Equal(t, http.StatusUnauthorized, w.Code)

			var resp base.Response
			json.Unmarshal(w.Body.Bytes(), &resp)
			require.False(t, resp.IsSuccess)
		})
	}
}
