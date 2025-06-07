package auth_test

import (
	"course_project/internal/config"
	"course_project/internal/services/auth"
	"github.com/stretchr/testify/assert"
	"testing"
)

func mockConfig() *config.Config {
	return &config.Config{
		AuthSecret:      "supersecretkey",
		TokenTTLMinutes: "60",
	}
}

func TestGenerateAndComparePasswordHash(t *testing.T) {
	service := auth.NewService(mockConfig())

	password := "mySecurePassword123"

	hash, err := service.GeneratePasswordHash(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	ok, err := service.CompareHashAndPassword(password, hash)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestCompareHashAndPassword_Invalid(t *testing.T) {
	service := auth.NewService(mockConfig())

	password := "correctPassword"
	wrongPassword := "wrongPassword"

	hash, _ := service.GeneratePasswordHash(password)

	ok, err := service.CompareHashAndPassword(wrongPassword, hash)
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestCreateAndVerifyAuthToken(t *testing.T) {
	service := auth.NewService(mockConfig())

	userID := "user-123"
	token, err := service.CreateAuthToken(userID)
	assert.NoError(t, err)
	assert.NotNil(t, token)

	subject, err := service.VerifyAuthToken(*token)
	assert.NoError(t, err)
	assert.Equal(t, userID, *subject)
}

func TestVerifyAuthToken_Invalid(t *testing.T) {
	service := auth.NewService(mockConfig())

	invalidToken := "invalid.token.value"

	subject, err := service.VerifyAuthToken(invalidToken)
	assert.Error(t, err)
	assert.Nil(t, subject)
}
