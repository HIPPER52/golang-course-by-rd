package client

import (
	"context"
	"course_project/internal/models"
)

type Repository interface {
	CountByPhone(ctx context.Context, phone string) (int64, error)
	Create(ctx context.Context, c *models.Client) error
}
