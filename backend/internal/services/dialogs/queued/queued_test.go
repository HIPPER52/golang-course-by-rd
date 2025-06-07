package queued_test

import (
	"context"
	"course_project/internal/models"
	"course_project/internal/services/dialogs/queued"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Add(ctx context.Context, dialog *models.QueuedDialog) error {
	args := m.Called(ctx, dialog)
	return args.Error(0)
}

func (m *MockRepo) FindByID(ctx context.Context, id string) (*models.QueuedDialog, error) {
	args := m.Called(ctx, id)
	val := args.Get(0)
	if val == nil {
		return nil, args.Error(1)
	}
	return val.(*models.QueuedDialog), args.Error(1)
}

func (m *MockRepo) ListAll(ctx context.Context) ([]models.QueuedDialog, error) {
	args := m.Called(ctx)
	val := args.Get(0)
	if val == nil {
		return nil, args.Error(1)
	}
	return val.([]models.QueuedDialog), args.Error(1)
}

func TestAdd_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	service := queued.NewService(mockRepo)

	dialog := &models.QueuedDialog{
		DialogBase: models.DialogBase{
			ID:         "dlg1",
			ClientID:   "client1",
			ClientName: "John",
		},
	}

	mockRepo.On("Add", ctx, dialog).Return(nil)

	err := service.Add(ctx, dialog)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAdd_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	service := queued.NewService(mockRepo)

	dialog := &models.QueuedDialog{
		DialogBase: models.DialogBase{ID: "dlg2"},
	}

	mockRepo.On("Add", ctx, dialog).Return(errors.New("db error"))

	err := service.Add(ctx, dialog)
	assert.Error(t, err)
}

func TestFindByID_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	service := queued.NewService(mockRepo)

	dialog := &models.QueuedDialog{
		DialogBase: models.DialogBase{ID: "dlg3"},
	}

	mockRepo.On("FindByID", ctx, "dlg3").Return(dialog, nil)

	result, err := service.FindByID(ctx, "dlg3")
	assert.NoError(t, err)
	assert.Equal(t, "dlg3", result.ID)
}

func TestFindByID_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	service := queued.NewService(mockRepo)

	mockRepo.On("FindByID", ctx, "dlg4").Return(nil, errors.New("not found"))

	result, err := service.FindByID(ctx, "dlg4")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestListAll_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	service := queued.NewService(mockRepo)

	dialogs := []models.QueuedDialog{
		{DialogBase: models.DialogBase{ID: "dlg5"}},
		{DialogBase: models.DialogBase{ID: "dlg6"}},
	}

	mockRepo.On("ListAll", ctx).Return(dialogs, nil)

	result, err := service.ListAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestListAll_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	service := queued.NewService(mockRepo)

	mockRepo.On("ListAll", ctx).Return(nil, errors.New("db fail"))

	result, err := service.ListAll(ctx)
	assert.Error(t, err)
	assert.Nil(t, result)
}
