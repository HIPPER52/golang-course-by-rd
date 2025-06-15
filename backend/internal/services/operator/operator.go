package operator

import (
	"context"
	"course_project/internal/dto"
	"course_project/internal/models"
	"course_project/internal/repository"
	repo "course_project/internal/repository/operator"
	"course_project/internal/services/logger"
	"errors"
	"fmt"
)

type Service struct {
	repo repo.Repository
}

func NewService(repo repo.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddOperator(ctx context.Context, dto dto.CreateOperatorDTO) (*models.Operator, error) {
	logger.Info(ctx, "Adding new operator: "+dto.Email)

	op, err := s.repo.AddOperator(ctx, dto)
	if err != nil {
		if errors.Is(err, repository.ErrOperatorAlreadyExists) {
			logger.Info(ctx, "Operator already exists: "+dto.Email)
		} else {
			logger.Error(ctx, fmt.Errorf("failed to add operator: %w", err))
		}
		return nil, err
	}

	logger.Info(ctx, "Operator successfully added: "+op.ID)
	return op, nil
}

func (s *Service) GetOperatorByEmail(ctx context.Context, email string) (*models.Operator, error) {
	logger.Info(ctx, "Fetching operator by email: "+email)

	op, err := s.repo.GetOperatorByEmail(ctx, email)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to get operator by email: %w", err))
		return nil, err
	}
	return op, nil
}

func (s *Service) GetOperatorByID(ctx context.Context, id string) (*models.Operator, error) {
	logger.Info(ctx, "Fetching operator by ID: "+id)

	op, err := s.repo.GetOperatorByID(ctx, id)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to get operator by ID: %w", err))
		return nil, err
	}
	return op, nil
}

func (s *Service) GetAllOperators(ctx context.Context) ([]*models.Operator, error) {
	logger.Info(ctx, "Fetching all operators")

	ops, err := s.repo.GetAllOperators(ctx)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to get all operators: %w", err))
		return nil, err
	}
	return ops, nil
}
