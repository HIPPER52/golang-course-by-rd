package queued

import (
	"context"
	"course_project/internal/models"
)

type Repository interface {
	Add(ctx context.Context, dialog *models.QueuedDialog) error
	FindByID(ctx context.Context, id string) (*models.QueuedDialog, error)
	ListAll(ctx context.Context) ([]models.QueuedDialog, error)
}
