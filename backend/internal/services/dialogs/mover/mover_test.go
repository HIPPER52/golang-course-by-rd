package mover_test

import (
	"context"
	"course_project/internal/models"
	"course_project/internal/services/dialogs/mover"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) FindQueuedByID(ctx context.Context, id string) (*models.QueuedDialog, error) {
	args := m.Called(ctx, id)
	val := args.Get(0)
	if val == nil {
		return nil, args.Error(1)
	}
	return val.(*models.QueuedDialog), args.Error(1)
}

func (m *MockRepo) InsertActive(ctx context.Context, dialog models.ActiveDialog) error {
	args := m.Called(ctx, dialog)
	return args.Error(0)
}

func (m *MockRepo) DeleteQueuedByID(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepo) FindActiveByID(ctx context.Context, id string) (*models.ActiveDialog, error) {
	args := m.Called(ctx, id)
	val := args.Get(0)
	if val == nil {
		return nil, args.Error(1)
	}
	return val.(*models.ActiveDialog), args.Error(1)
}

func (m *MockRepo) InsertArchived(ctx context.Context, dialog models.ArchivedDialog) error {
	args := m.Called(ctx, dialog)
	return args.Error(0)
}

func (m *MockRepo) DeleteActiveByID(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTakeDialog_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	service := mover.NewService(mockRepo)

	queued := &models.QueuedDialog{
		DialogBase: models.DialogBase{
			ID:         "dialog123",
			ClientID:   "client123",
			ClientName: "John",
			StartedAt:  time.Now().Add(-10 * time.Minute).UTC(),
		},
	}

	mockRepo.On("FindQueuedByID", ctx, "dialog123").Return(queued, nil)
	mockRepo.On("InsertActive", ctx, mock.AnythingOfType("models.ActiveDialog")).Return(nil)
	mockRepo.On("DeleteQueuedByID", ctx, "dialog123").Return(nil)

	err := service.TakeDialog(ctx, "dialog123", "operator123")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTakeDialog_QueuedNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	service := mover.NewService(mockRepo)

	mockRepo.On("FindQueuedByID", ctx, "dialog123").Return(nil, nil)

	err := service.TakeDialog(ctx, "dialog123", "operator123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestCloseDialog_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	service := mover.NewService(mockRepo)

	active := &models.ActiveDialog{
		DialogBase: models.DialogBase{
			ID:         "dialog456",
			ClientID:   "client456",
			ClientName: "Alice",
			OperatorID: "op456",
			StartedAt:  time.Now().Add(-15 * time.Minute).UTC(),
		},
	}

	mockRepo.On("FindActiveByID", ctx, "dialog456").Return(active, nil)
	mockRepo.On("InsertArchived", ctx, mock.AnythingOfType("models.ArchivedDialog")).Return(nil)
	mockRepo.On("DeleteActiveByID", ctx, "dialog456").Return(nil)

	err := service.CloseDialog(ctx, "dialog456")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCloseDialog_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	service := mover.NewService(mockRepo)

	mockRepo.On("FindActiveByID", ctx, "dialog456").Return(nil, nil)

	err := service.CloseDialog(ctx, "dialog456")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
