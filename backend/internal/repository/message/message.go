package message

import (
	"context"
	"course_project/internal/models"
)

type Repository interface {
	FindByRoomID(ctx context.Context, roomID string) ([]models.Message, error)
}

type DialogFinder interface {
	FindDialogByID(ctx context.Context, id string) (*models.ArchivedDialog, error)
}
