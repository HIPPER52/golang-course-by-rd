package auth

import (
	"context"
	"course_project/internal/constants/roles"
	"course_project/internal/services/logger"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	UserID string     `json:"sub"`
	Role   roles.Role `json:"role"`
	jwt.RegisteredClaims
}

func (s *Service) CreateAuthToken(userId string, role roles.Role) (*string, error) {
	utcNow := time.Now().UTC()

	claims := AuthClaims{
		UserID: userId,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			IssuedAt:  jwt.NewNumericDate(utcNow),
			ExpiresAt: jwt.NewNumericDate(utcNow.Add(s.tokenTTLMinutes)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.cfg.AuthSecret))
	if err != nil {
		logger.Error(context.Background(), fmt.Errorf("failed to sign token"))
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	logger.Info(context.Background(), "Token generated for userId: "+userId)
	return &signedToken, nil
}

func (s *Service) VerifyAuthToken(tokenStr string) (*AuthClaims, error) {
	claims := &AuthClaims{}

	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.AuthSecret), nil
	})

	if err != nil {
		logger.Error(context.Background(), fmt.Errorf("failed to verify token"))
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return claims, nil
}
