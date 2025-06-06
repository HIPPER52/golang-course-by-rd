package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func (s *Service) CreateAuthToken(userId string) (*string, error) {
	utcNow := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Subject:   userId,
		IssuedAt:  jwt.NewNumericDate(utcNow),
		ExpiresAt: jwt.NewNumericDate(utcNow.Add(s.tokenTTLMinutes)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.cfg.AuthSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &signedToken, nil
}

func (s *Service) VerifyAuthToken(token string) (*string, error) {
	claims := jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.AuthSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return &claims.Subject, nil
}
