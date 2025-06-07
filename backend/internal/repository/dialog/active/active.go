package active

import (
	"context"
	"course_project/internal/models"
)

type Repository interface {
	Add(ctx context.Context, dialog *models.ActiveDialog) error
	ListAll(ctx context.Context) ([]models.ActiveDialog, error)
	FindByID(ctx context.Context, id string) (*models.ActiveDialog, error)
	FindByOperatorID(ctx context.Context, operatorID string) ([]models.ActiveDialog, error)
	CountByOperator(ctx context.Context, operatorID string) (int, error)
}
