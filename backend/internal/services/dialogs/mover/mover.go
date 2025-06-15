package mover

import (
	"context"
	"course_project/internal/models"
	"course_project/internal/repository/dialog/mover"
	"course_project/internal/services/logger"
	"fmt"
	"time"
)

type Service struct {
	repo mover.Repository
}

func NewService(repo mover.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) TakeDialog(ctx context.Context, dialogID, operatorID string) error {
	logger.Info(ctx, "Operator "+operatorID+" is taking dialog: "+dialogID)

	queued, err := s.repo.FindQueuedByID(ctx, dialogID)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to find queued dialog: %w", err))
		return err
	}
	if queued == nil {
		logger.Error(ctx, fmt.Errorf("dialog %s not found in queue", dialogID))
		return fmt.Errorf("dialog %q not found in queue", dialogID)
	}

	active := models.ActiveDialog{
		DialogBase: models.DialogBase{
			ID:            queued.ID,
			ClientID:      queued.ClientID,
			ClientName:    queued.ClientName,
			ClientPhone:   queued.ClientPhone,
			ClientIP:      queued.ClientIP,
			OperatorID:    operatorID,
			StartedAt:     time.Now().UTC(),
			LastMessageAt: queued.StartedAt,
		},
	}

	if err := s.repo.InsertActive(ctx, active); err != nil {
		logger.Error(ctx, fmt.Errorf("failed to insert active dialog: %w", err))
		return err
	}

	if err := s.repo.DeleteQueuedByID(ctx, dialogID); err != nil {
		logger.Error(ctx, fmt.Errorf("failed to delete queued dialog: %w", err))
		return err
	}

	logger.Info(ctx, "Dialog "+dialogID+" successfully taken by operator "+operatorID)
	return nil
}

func (s *Service) CloseDialog(ctx context.Context, dialogID string) error {
	logger.Info(ctx, "Closing dialog: "+dialogID)

	active, err := s.repo.FindActiveByID(ctx, dialogID)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to find active dialog: %w", err))
		return err
	}
	if active == nil {
		logger.Error(ctx, fmt.Errorf("dialog %s not found in active", dialogID))
		return fmt.Errorf("dialog %q not found in active", dialogID)
	}

	archived := models.ArchivedDialog{
		DialogBase: models.DialogBase{
			ID:            active.ID,
			ClientID:      active.ClientID,
			ClientName:    active.ClientName,
			ClientPhone:   active.ClientPhone,
			ClientIP:      active.ClientIP,
			OperatorID:    active.OperatorID,
			StartedAt:     active.StartedAt,
			LastMessageAt: active.LastMessageAt,
			EndedAt:       time.Now().UTC(),
		},
	}

	if err := s.repo.InsertArchived(ctx, archived); err != nil {
		logger.Error(ctx, fmt.Errorf("failed to insert archived dialog: %w", err))
		return err
	}

	if err := s.repo.DeleteActiveByID(ctx, dialogID); err != nil {
		logger.Error(ctx, fmt.Errorf("failed to delete active dialog: %w", err))
		return err
	}

	logger.Info(ctx, "Dialog "+dialogID+" successfully archived and closed")
	return nil
}
