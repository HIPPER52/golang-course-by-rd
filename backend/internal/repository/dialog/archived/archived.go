package archived

import (
	"context"
	"course_project/internal/models"
)

type Repository interface {
	Add(ctx context.Context, dialog *models.ArchivedDialog) error
	FindByOperator(ctx context.Context, operatorID string) ([]models.ArchivedDialog, error)
	CountByOperator(ctx context.Context, operatorID string) (int, error)
}
