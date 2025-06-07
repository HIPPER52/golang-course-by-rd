package message_test

import (
	"context"
	"course_project/internal/models"
	"course_project/internal/services/message"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) FindByRoomID(ctx context.Context, roomID string) ([]models.Message, error) {
	args := m.Called(ctx, roomID)
	return args.Get(0).([]models.Message), args.Error(1)
}

type MockDialogFinder struct {
	mock.Mock
}

func (m *MockDialogFinder) FindDialogByID(ctx context.Context, id string) (*models.ArchivedDialog, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.ArchivedDialog), args.Error(1)
}

func TestFindByRoomID_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	mockFinder := new(MockDialogFinder)
	svc := message.NewService(mockRepo, mockFinder)

	roomID := "room123"
	clientID := "client123"
	dialog := &models.ArchivedDialog{
		DialogBase: models.DialogBase{
			ID:       roomID,
			ClientID: clientID,
		},
	}
	msgs := []models.Message{{ID: "msg1"}, {ID: "msg2"}}

	mockFinder.On("FindDialogByID", ctx, roomID).Return(dialog, nil)
	mockRepo.On("FindByRoomID", ctx, roomID).Return(msgs, nil)

	result, err := svc.FindByRoomID(ctx, roomID, clientID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockFinder.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestFindByRoomID_DialogNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	mockFinder := new(MockDialogFinder)
	svc := message.NewService(mockRepo, mockFinder)

	roomID := "room404"
	mockFinder.On("FindDialogByID", ctx, roomID).Return((*models.ArchivedDialog)(nil), errors.New("not found"))

	result, err := svc.FindByRoomID(ctx, roomID, "client123")

	assert.Nil(t, result)
	assert.EqualError(t, err, "not found")
	mockFinder.AssertExpectations(t)
}

func TestFindByRoomID_AccessDenied(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	mockFinder := new(MockDialogFinder)
	svc := message.NewService(mockRepo, mockFinder)

	roomID := "room123"
	dialog := &models.ArchivedDialog{
		DialogBase: models.DialogBase{
			ID:       roomID,
			ClientID: "other_client",
		},
	}

	mockFinder.On("FindDialogByID", ctx, roomID).Return(dialog, nil)

	result, err := svc.FindByRoomID(ctx, roomID, "client123")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, message.ErrAccessDenied)
	mockFinder.AssertExpectations(t)
}

func TestFindByRoomID_RepoError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepo)
	mockFinder := new(MockDialogFinder)
	svc := message.NewService(mockRepo, mockFinder)

	roomID := "room123"
	clientID := "client123"
	dialog := &models.ArchivedDialog{
		DialogBase: models.DialogBase{
			ID:       roomID,
			ClientID: clientID,
		},
	}

	mockFinder.On("FindDialogByID", ctx, roomID).Return(dialog, nil)
	mockRepo.On("FindByRoomID", ctx, roomID).Return([]models.Message(nil), errors.New("db error"))

	result, err := svc.FindByRoomID(ctx, roomID, clientID)

	assert.Nil(t, result)
	assert.EqualError(t, err, "db error")
	mockFinder.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
