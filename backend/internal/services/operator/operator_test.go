package operator_test

import (
	"context"
	"course_project/internal/dto"
	"course_project/internal/models"
	"course_project/internal/repository"
	"course_project/internal/services/operator"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockOperatorRepo struct {
	mock.Mock
}

func (m *MockOperatorRepo) AddOperator(ctx context.Context, dto dto.CreateOperatorDTO) (*models.Operator, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(*models.Operator), args.Error(1)
}

func (m *MockOperatorRepo) GetOperatorByEmail(ctx context.Context, email string) (*models.Operator, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*models.Operator), args.Error(1)
}

func (m *MockOperatorRepo) GetOperatorByID(ctx context.Context, id string) (*models.Operator, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Operator), args.Error(1)
}

func (m *MockOperatorRepo) GetAllOperators(ctx context.Context) ([]*models.Operator, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Operator), args.Error(1)
}

func TestAddOperator_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockOperatorRepo)
	service := operator.NewService(mockRepo)

	req := dto.CreateOperatorDTO{
		Username: "John",
		Email:    "john@example.com",
		Role:     "operator",
	}

	expected := &models.Operator{ID: "op123", Email: req.Email}

	mockRepo.On("AddOperator", ctx, req).Return(expected, nil)

	result, err := service.AddOperator(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestAddOperator_AlreadyExists(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockOperatorRepo)
	service := operator.NewService(mockRepo)

	req := dto.CreateOperatorDTO{Email: "john@example.com"}

	mockRepo.On("AddOperator", ctx, req).Return((*models.Operator)(nil), repository.ErrOperatorAlreadyExists)

	result, err := service.AddOperator(ctx, req)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, repository.ErrOperatorAlreadyExists)
	mockRepo.AssertExpectations(t)
}

func TestAddOperator_UnexpectedError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockOperatorRepo)
	service := operator.NewService(mockRepo)

	req := dto.CreateOperatorDTO{Email: "fail@example.com"}

	mockRepo.On("AddOperator", ctx, req).Return((*models.Operator)(nil), errors.New("db error"))

	result, err := service.AddOperator(ctx, req)

	assert.Nil(t, result)
	assert.EqualError(t, err, "db error")
	mockRepo.AssertExpectations(t)
}

func TestGetOperatorByEmail_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockOperatorRepo)
	service := operator.NewService(mockRepo)

	email := "admin@example.com"
	expected := &models.Operator{ID: "id1", Email: email}

	mockRepo.On("GetOperatorByEmail", ctx, email).Return(expected, nil)

	result, err := service.GetOperatorByEmail(ctx, email)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetOperatorByID_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockOperatorRepo)
	service := operator.NewService(mockRepo)

	id := "op123"
	expected := &models.Operator{ID: id, Email: "admin@example.com"}

	mockRepo.On("GetOperatorByID", ctx, id).Return(expected, nil)

	result, err := service.GetOperatorByID(ctx, id)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetAllOperators_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockOperatorRepo)
	service := operator.NewService(mockRepo)

	expected := []*models.Operator{
		{ID: "op1", Email: "a@example.com"},
		{ID: "op2", Email: "b@example.com"},
	}

	mockRepo.On("GetAllOperators", ctx).Return(expected, nil)

	result, err := service.GetAllOperators(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}
