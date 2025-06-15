package message

import (
	"context"
	"course_project/internal/models"
	repo "course_project/internal/repository/message"
	"course_project/internal/services/logger"
	"errors"
	"fmt"
)

var (
	ErrAccessDenied   = errors.New("you don't have access to this dialog")
	ErrDialogNotFound = errors.New("dialog not found")
)

type Service struct {
	repo         repo.Repository
	dialogFinder repo.DialogFinder
}

func NewService(repo repo.Repository, finder repo.DialogFinder) *Service {
	return &Service{
		repo:         repo,
		dialogFinder: finder,
	}
}

func (s *Service) FindByRoomID(ctx context.Context, roomID string, clientID string) ([]models.Message, error) {
	logger.Info(ctx, "Searching dialog and messages for room: "+roomID)

	dialog, err := s.dialogFinder.FindDialogByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	if clientID != "" && dialog.ClientID != clientID {
		logger.Info(ctx, "Access denied to room "+roomID+" for client "+clientID)
		return nil, ErrAccessDenied
	}

	messages, err := s.repo.FindByRoomID(ctx, roomID)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to find messages: %w", err))
		return nil, err
	}

	logger.Info(ctx, fmt.Sprintf("Found %d messages in room %s", len(messages), roomID))
	return messages, nil
}
