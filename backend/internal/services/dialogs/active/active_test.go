package active_test

import (
	"context"
	"course_project/internal/models"
	"course_project/internal/services/dialogs/active"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Add(ctx context.Context, dialog *models.ActiveDialog) error {
	args := m.Called(ctx, dialog)
	return args.Error(0)
}

func (m *MockRepo) ListAll(ctx context.Context) ([]models.ActiveDialog, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.ActiveDialog), args.Error(1)
}

func (m *MockRepo) FindByID(ctx context.Context, id string) (*models.ActiveDialog, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.ActiveDialog), args.Error(1)
}

func (m *MockRepo) FindByOperatorID(ctx context.Context, operatorID string) ([]models.ActiveDialog, error) {
	args := m.Called(ctx, operatorID)
	return args.Get(0).([]models.ActiveDialog), args.Error(1)
}

func (m *MockRepo) CountByOperator(ctx context.Context, operatorID string) (int, error) {
	args := m.Called(ctx, operatorID)
	return args.Int(0), args.Error(1)
}

func TestAdd_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	svc := active.NewService(mockRepo)

	dialog := &models.ActiveDialog{DialogBase: models.DialogBase{ID: "dialog123"}}
	mockRepo.On("Add", ctx, dialog).Return(nil)

	err := svc.Add(ctx, dialog)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAdd_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	svc := active.NewService(mockRepo)

	dialog := &models.ActiveDialog{DialogBase: models.DialogBase{ID: "dialog123"}}
	mockRepo.On("Add", ctx, dialog).Return(errors.New("db error"))

	err := svc.Add(ctx, dialog)
	assert.EqualError(t, err, "db error")
	mockRepo.AssertExpectations(t)
}

func TestListAll_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	svc := active.NewService(mockRepo)

	expected := []models.ActiveDialog{{}, {}}
	mockRepo.On("ListAll", ctx).Return(expected, nil)

	result, err := svc.ListAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestListAll_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	svc := active.NewService(mockRepo)

	mockRepo.On("ListAll", ctx).Return([]models.ActiveDialog(nil), errors.New("db error"))

	result, err := svc.ListAll(ctx)
	assert.Nil(t, result)
	assert.EqualError(t, err, "db error")
	mockRepo.AssertExpectations(t)
}

func TestFindByID_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	svc := active.NewService(mockRepo)

	dialog := &models.ActiveDialog{DialogBase: models.DialogBase{ID: "dialog123"}}
	mockRepo.On("FindByID", ctx, "dialog123").Return(dialog, nil)

	result, err := svc.FindByID(ctx, "dialog123")
	assert.NoError(t, err)
	assert.Equal(t, dialog, result)
	mockRepo.AssertExpectations(t)
}

func TestFindByID_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	svc := active.NewService(mockRepo)

	mockRepo.On("FindByID", ctx, "dialog123").Return((*models.ActiveDialog)(nil), errors.New("not found"))

	result, err := svc.FindByID(ctx, "dialog123")
	assert.Nil(t, result)
	assert.EqualError(t, err, "not found")
	mockRepo.AssertExpectations(t)
}

func TestFindByOperatorID_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	svc := active.NewService(mockRepo)

	dialogs := []models.ActiveDialog{{}, {}}
	mockRepo.On("FindByOperatorID", ctx, "op123").Return(dialogs, nil)

	result, err := svc.FindByOperatorID(ctx, "op123")
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestFindByOperatorID_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	svc := active.NewService(mockRepo)

	mockRepo.On("FindByOperatorID", ctx, "op123").Return([]models.ActiveDialog(nil), errors.New("db error"))

	result, err := svc.FindByOperatorID(ctx, "op123")
	assert.Nil(t, result)
	assert.EqualError(t, err, "db error")
	mockRepo.AssertExpectations(t)
}

func TestCountByOperator_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	svc := active.NewService(mockRepo)

	mockRepo.On("CountByOperator", ctx, "op123").Return(5, nil)

	count, err := svc.CountByOperator(ctx, "op123")
	assert.NoError(t, err)
	assert.Equal(t, 5, count)
	mockRepo.AssertExpectations(t)
}

func TestCountByOperator_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	svc := active.NewService(mockRepo)

	mockRepo.On("CountByOperator", ctx, "op123").Return(0, errors.New("db error"))

	count, err := svc.CountByOperator(ctx, "op123")
	assert.Equal(t, 0, count)
	assert.EqualError(t, err, "db error")
	mockRepo.AssertExpectations(t)
}
