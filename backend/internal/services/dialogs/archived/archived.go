package archived

import (
	"context"
	"course_project/internal/models"
	"course_project/internal/repository"
	repo "course_project/internal/repository/dialog/archived"
	"course_project/internal/services/logger"
	"errors"
	"fmt"
)

var (
	ErrArchivedDialogExists = repository.ErrArchivedDialogExists
)

type Service struct {
	repo repo.Repository
}

func NewService(repo repo.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Add(ctx context.Context, dialog *models.ArchivedDialog) error {
	logger.Info(ctx, "Archiving dialog: "+dialog.ID)

	err := s.repo.Add(ctx, dialog)
	if errors.Is(err, ErrArchivedDialogExists) {
		logger.Error(ctx, fmt.Errorf("dialog already archived: %w", err))
		return ErrArchivedDialogExists
	}
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to archive dialog: %w", err))
	}
	return err
}

func (s *Service) FindByOperator(ctx context.Context, operatorID string) ([]models.ArchivedDialog, error) {
	logger.Info(ctx, "Finding archived dialogs by operator: "+operatorID)

	dialogs, err := s.repo.FindByOperator(ctx, operatorID)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to find archived dialogs: %w", err))
	}
	return dialogs, err
}

func (s *Service) CountByOperator(ctx context.Context, operatorID string) (int, error) {
	logger.Info(ctx, "Counting archived dialogs by operator: "+operatorID)

	count, err := s.repo.CountByOperator(ctx, operatorID)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to count archived dialogs: %w", err))
	}
	return count, err
}
