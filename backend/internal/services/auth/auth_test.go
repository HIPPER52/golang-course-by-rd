package auth_test

import (
	"course_project/internal/config"
	"course_project/internal/constants/roles"
	"course_project/internal/services/auth"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockConfig() *config.Config {
	return &config.Config{
		AuthSecret:      "supersecretkey",
		TokenTTLMinutes: "60",
	}
}

func TestCreateAndVerifyAuthToken(t *testing.T) {
	service := auth.NewService(mockConfig())

	userID := "user-123"
	userRole := roles.Operator
	token, err := service.CreateAuthToken(userID, userRole)
	assert.NoError(t, err)
	assert.NotNil(t, token)

	subject, err := service.VerifyAuthToken(*token)
	assert.NoError(t, err)
	assert.Equal(t, userID, subject.UserID)
}

func TestVerifyAuthToken_Invalid(t *testing.T) {
	service := auth.NewService(mockConfig())

	invalidToken := "invalid.token.value"

	subject, err := service.VerifyAuthToken(invalidToken)
	assert.Error(t, err)
	assert.Nil(t, subject)
}
