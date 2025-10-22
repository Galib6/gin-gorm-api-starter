package service_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zetsux/gin-gorm-api-starter/core/service"
)

func TestJWTService_GenerateValidateAndGetAttrs(t *testing.T) {
	jwtS := service.NewJWTService()

	token := jwtS.GenerateToken("user-123", "user")
	require.NotEmpty(t, token)

	parsed, err := jwtS.ValidateToken(token)
	require.NoError(t, err)
	require.True(t, parsed.Valid)

	id, role, err := jwtS.GetAttrByToken(token)
	require.NoError(t, err)
	require.Equal(t, "user-123", id)
	require.Equal(t, "user", role)
}
