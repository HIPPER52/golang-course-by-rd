package operator

import (
	"context"
	"course_project/internal/dto"
	"course_project/internal/models"
)

type Repository interface {
	AddOperator(ctx context.Context, dto dto.CreateOperatorDTO) (*models.Operator, error)
	GetOperatorByEmail(ctx context.Context, email string) (*models.Operator, error)
	GetOperatorByID(ctx context.Context, id string) (*models.Operator, error)
	GetAllOperators(ctx context.Context) ([]*models.Operator, error)
}
