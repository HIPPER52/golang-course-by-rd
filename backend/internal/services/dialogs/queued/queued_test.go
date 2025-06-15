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

func TestAdd(t *testing.T) {
	ctx := context.Background()
	dialog := &models.QueuedDialog{DialogBase: models.DialogBase{ID: "dlg1"}}

	tests := []struct {
		name    string
		setup   func(m *MockRepo)
		wantErr bool
	}{
		{
			name: "success",
			setup: func(m *MockRepo) {
				m.On("Add", ctx, dialog).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "db error",
			setup: func(m *MockRepo) {
				m.On("Add", ctx, dialog).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			tt.setup(mockRepo)
			service := queued.NewService(mockRepo)

			err := service.Add(ctx, dialog)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestFindByID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		inputID  string
		mockResp *models.QueuedDialog
		mockErr  error
		wantNil  bool
	}{
		{
			name:     "success",
			inputID:  "dlg3",
			mockResp: &models.QueuedDialog{DialogBase: models.DialogBase{ID: "dlg3"}},
			mockErr:  nil,
			wantNil:  false,
		},
		{
			name:     "not found",
			inputID:  "dlg4",
			mockResp: nil,
			mockErr:  errors.New("not found"),
			wantNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			mockRepo.On("FindByID", ctx, tt.inputID).Return(tt.mockResp, tt.mockErr)
			service := queued.NewService(mockRepo)

			result, err := service.FindByID(ctx, tt.inputID)
			if tt.wantNil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tt.mockResp.ID, result.ID)
			}
			assert.Equal(t, tt.mockErr != nil, err != nil)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestListAll(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		mockResp []models.QueuedDialog
		mockErr  error
		wantErr  bool
	}{
		{
			name:     "success",
			mockResp: []models.QueuedDialog{{DialogBase: models.DialogBase{ID: "dlg5"}}, {DialogBase: models.DialogBase{ID: "dlg6"}}},
			mockErr:  nil,
			wantErr:  false,
		},
		{
			name:     "db fail",
			mockResp: nil,
			mockErr:  errors.New("db fail"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			mockRepo.On("ListAll", ctx).Return(tt.mockResp, tt.mockErr)
			service := queued.NewService(mockRepo)

			result, err := service.ListAll(ctx)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, len(tt.mockResp))
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
