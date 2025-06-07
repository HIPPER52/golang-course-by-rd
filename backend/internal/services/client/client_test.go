package client_test

import (
	"context"
	"course_project/internal/dto"
	"course_project/internal/models"
	"course_project/internal/repository"
	"course_project/internal/services/client"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockClientRepo struct {
	mock.Mock
}

func (m *MockClientRepo) CountByPhone(ctx context.Context, phone string) (int64, error) {
	args := m.Called(ctx, phone)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockClientRepo) Create(ctx context.Context, c *models.Client) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func TestRegisterClient_Success(t *testing.T) {
	mockRepo := new(MockClientRepo)
	service := client.NewService(mockRepo)

	ctx := context.Background()
	req := dto.RegisterClientDTO{
		Name:  "John Doe",
		Phone: "+380931234567",
	}

	mockRepo.On("CountByPhone", ctx, req.Phone).Return(int64(0), nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Client")).Return(nil)

	cl, err := service.RegisterClient(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, cl)
	assert.Equal(t, req.Name, cl.Name)
	assert.Equal(t, req.Phone, cl.Phone)
	assert.WithinDuration(t, time.Now(), cl.CreatedAt, time.Second)

	mockRepo.AssertExpectations(t)
}

func TestRegisterClient_AlreadyExists(t *testing.T) {
	mockRepo := new(MockClientRepo)
	service := client.NewService(mockRepo)

	ctx := context.Background()
	req := dto.RegisterClientDTO{
		Name:  "Jane",
		Phone: "+380931234567",
	}

	mockRepo.On("CountByPhone", ctx, req.Phone).Return(int64(1), nil)

	cl, err := service.RegisterClient(ctx, req)

	assert.Nil(t, cl)
	assert.ErrorIs(t, err, repository.ErrClientAlreadyExists)
	mockRepo.AssertExpectations(t)
}

func TestRegisterClient_CountError(t *testing.T) {
	mockRepo := new(MockClientRepo)
	service := client.NewService(mockRepo)

	ctx := context.Background()
	req := dto.RegisterClientDTO{
		Name:  "Error",
		Phone: "+380931234567",
	}

	mockRepo.On("CountByPhone", ctx, req.Phone).Return(int64(0), errors.New("db error"))

	cl, err := service.RegisterClient(ctx, req)

	assert.Nil(t, cl)
	assert.Contains(t, err.Error(), "failed to check existing clients")
	mockRepo.AssertExpectations(t)
}

func TestRegisterClient_CreateError(t *testing.T) {
	mockRepo := new(MockClientRepo)
	service := client.NewService(mockRepo)

	ctx := context.Background()
	req := dto.RegisterClientDTO{
		Name:  "New",
		Phone: "+380931234567",
	}

	mockRepo.On("CountByPhone", ctx, req.Phone).Return(int64(0), nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Client")).Return(errors.New("insert error"))

	cl, err := service.RegisterClient(ctx, req)

	assert.Nil(t, cl)
	assert.Contains(t, err.Error(), "failed to insert client")
	mockRepo.AssertExpectations(t)
}
