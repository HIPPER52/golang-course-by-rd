package queued

import (
	"context"
	"course_project/internal/models"
	repo "course_project/internal/repository/dialog/queued"
	"course_project/internal/services/logger"
	"fmt"
)

type Service struct {
	repo repo.Repository
}

func NewService(r repo.Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Add(ctx context.Context, dialog *models.QueuedDialog) error {
	logger.Info(ctx, "Adding queued dialog: "+dialog.ID)
	err := s.repo.Add(ctx, dialog)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to add queued dialog: %w", err))
	}
	return err
}

func (s *Service) FindByID(ctx context.Context, id string) (*models.QueuedDialog, error) {
	logger.Info(ctx, "Finding queued dialog by ID: "+id)
	d, err := s.repo.FindByID(ctx, id)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to find queued dialog by ID: %w", err))
	}
	return d, err
}

func (s *Service) ListAll(ctx context.Context) ([]models.QueuedDialog, error) {
	logger.Info(ctx, "Listing all queued dialogs")
	dialogs, err := s.repo.ListAll(ctx)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to list queued dialogs: %w", err))
	}
	return dialogs, err
}
