package auth

import (
	"context"
	"course_project/internal/config"
	"course_project/internal/services/logger"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type Service struct {
	cfg             *config.Config
	tokenTTLMinutes time.Duration
}

func NewService(cfg *config.Config) *Service {
	minutesInt, err := strconv.Atoi(cfg.TokenTTLMinutes)
	if err != nil {
		panic("config error, tokenTTLMinutes" + err.Error())
	}

	logger.Info(context.Background(), "Auth service initialized with token TTL: "+cfg.TokenTTLMinutes+" minutes")

	return &Service{
		tokenTTLMinutes: time.Duration(minutesInt) * time.Minute,
		cfg:             cfg,
	}
}

func (s *Service) GeneratePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *Service) CompareHashAndPassword(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		logger.Error(context.Background(), fmt.Errorf("password mismatch"))
		return false, nil
	}
	if err != nil {
		logger.Error(context.Background(), fmt.Errorf("error comparing password hash: %w", err))
		return false, err
	}
	return true, nil
}
