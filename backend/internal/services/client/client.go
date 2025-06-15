package client

import (
	"context"
	"course_project/internal/dto"
	"course_project/internal/models"
	"course_project/internal/repository"
	"course_project/internal/repository/client"
	"course_project/internal/services/logger"
	"fmt"
	"github.com/oklog/ulid/v2"
	"time"
)

var ErrClientAlreadyExists = repository.ErrClientAlreadyExists

type Service struct {
	repo client.Repository
}

func NewService(repo client.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterClient(ctx context.Context, dto dto.RegisterClientDTO) (*models.Client, error) {
	t := time.Now().UTC()

	logger.Info(ctx, "Registering new client: "+dto.Phone)

	count, err := s.repo.CountByPhone(ctx, dto.Phone)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to check existing clients"))
		return nil, fmt.Errorf("failed to check existing clients: %w", err)
	}
	if count > 0 {
		logger.Info(ctx, "Client already exists: "+dto.Phone)
		return nil, ErrClientAlreadyExists
	}

	client := &models.Client{
		ID:        ulid.Make().String(),
		Name:      dto.Name,
		Phone:     dto.Phone,
		CreatedAt: t,
	}

	if err := s.repo.Create(ctx, client); err != nil {
		logger.Error(ctx, fmt.Errorf("failed to insert client"))
		return nil, fmt.Errorf("failed to insert client: %w", err)
	}

	logger.Info(ctx, "New client registered: "+client.ID)
	return client, nil
}
