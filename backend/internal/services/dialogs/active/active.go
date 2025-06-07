package active

import (
	"context"
	"course_project/internal/models"
	repo "course_project/internal/repository/dialog/active"
	"course_project/internal/services/logger"
	"fmt"
)

type Service struct {
	repo repo.Repository
}

func NewService(repo repo.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Add(ctx context.Context, dialog *models.ActiveDialog) error {
	logger.Info(ctx, "Adding active dialog: "+dialog.ID)
	err := s.repo.Add(ctx, dialog)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to add active dialog: %w", err))
	}
	return err
}

func (s *Service) ListAll(ctx context.Context) ([]models.ActiveDialog, error) {
	logger.Info(ctx, "Listing all active dialogs")
	dialogs, err := s.repo.ListAll(ctx)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to list active dialogs: %w", err))
	}
	return dialogs, err
}

func (s *Service) FindByID(ctx context.Context, id string) (*models.ActiveDialog, error) {
	logger.Info(ctx, "Finding active dialog by ID: "+id)
	dialog, err := s.repo.FindByID(ctx, id)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to find active dialog by ID: %w", err))
	}
	return dialog, err
}

func (s *Service) FindByOperatorID(ctx context.Context, operatorID string) ([]models.ActiveDialog, error) {
	logger.Info(ctx, "Finding active dialogs by operator: "+operatorID)
	dialogs, err := s.repo.FindByOperatorID(ctx, operatorID)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to find active dialogs by operator: %w", err))
	}
	return dialogs, err
}

func (s *Service) CountByOperator(ctx context.Context, operatorID string) (int, error) {
	logger.Info(ctx, "Counting active dialogs by operator: "+operatorID)
	count, err := s.repo.CountByOperator(ctx, operatorID)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to count active dialogs by operator: %w", err))
	}
	return count, err
}
