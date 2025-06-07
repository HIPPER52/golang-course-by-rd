package mover

import (
	"context"
	"course_project/internal/models"
)

type Repository interface {
	FindQueuedByID(ctx context.Context, id string) (*models.QueuedDialog, error)
	InsertActive(ctx context.Context, dialog models.ActiveDialog) error
	DeleteQueuedByID(ctx context.Context, id string) error

	FindActiveByID(ctx context.Context, id string) (*models.ActiveDialog, error)
	InsertArchived(ctx context.Context, dialog models.ArchivedDialog) error
	DeleteActiveByID(ctx context.Context, id string) error
}
