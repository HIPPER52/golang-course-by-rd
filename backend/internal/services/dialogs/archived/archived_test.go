package archived_test

import (
	"context"
	"course_project/internal/models"
	"course_project/internal/repository"
	"course_project/internal/services/dialogs/archived"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockArchivedRepo struct {
	AddFunc             func(ctx context.Context, d *models.ArchivedDialog) error
	FindByOperatorFunc  func(ctx context.Context, opID string) ([]models.ArchivedDialog, error)
	CountByOperatorFunc func(ctx context.Context, opID string) (int, error)
}

func (m *MockArchivedRepo) Add(ctx context.Context, d *models.ArchivedDialog) error {
	return m.AddFunc(ctx, d)
}

func (m *MockArchivedRepo) FindByOperator(ctx context.Context, opID string) ([]models.ArchivedDialog, error) {
	return m.FindByOperatorFunc(ctx, opID)
}

func (m *MockArchivedRepo) CountByOperator(ctx context.Context, opID string) (int, error) {
	return m.CountByOperatorFunc(ctx, opID)
}

func TestAdd_Success(t *testing.T) {
	mockRepo := &MockArchivedRepo{
		AddFunc: func(ctx context.Context, d *models.ArchivedDialog) error {
			return nil
		},
	}

	service := archived.NewService(mockRepo)
	err := service.Add(context.TODO(), &models.ArchivedDialog{
		DialogBase: models.DialogBase{ID: "dialog123"},
	})
	assert.NoError(t, err)
}

func TestAdd_AlreadyExists(t *testing.T) {
	mockRepo := &MockArchivedRepo{
		AddFunc: func(ctx context.Context, d *models.ArchivedDialog) error {
			return repository.ErrArchivedDialogExists
		},
	}

	service := archived.NewService(mockRepo)
	err := service.Add(context.TODO(), &models.ArchivedDialog{
		DialogBase: models.DialogBase{ID: "dialog123"},
	})
	assert.ErrorIs(t, err, repository.ErrArchivedDialogExists)
}

func TestAdd_Error(t *testing.T) {
	mockRepo := &MockArchivedRepo{
		AddFunc: func(ctx context.Context, d *models.ArchivedDialog) error {
			return errors.New("db failure")
		},
	}

	service := archived.NewService(mockRepo)
	err := service.Add(context.TODO(), &models.ArchivedDialog{
		DialogBase: models.DialogBase{ID: "dialog123"},
	})
	assert.EqualError(t, err, "db failure")
}

func TestFindByOperator_Error(t *testing.T) {
	mockRepo := &MockArchivedRepo{
		FindByOperatorFunc: func(ctx context.Context, opID string) ([]models.ArchivedDialog, error) {
			return nil, errors.New("find error")
		},
	}

	service := archived.NewService(mockRepo)
	_, err := service.FindByOperator(context.TODO(), "op123")
	assert.EqualError(t, err, "find error")
}
