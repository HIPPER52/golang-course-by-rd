package operator

import (
	"context"
	"course_project/internal/clients"
	"course_project/internal/constants"
	"course_project/internal/dto"
	"course_project/internal/models"
	"course_project/internal/services/logger"
	"errors"
	"fmt"
	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var (
	ErrOperatorAlreadyExists = errors.New("operator already exists")
)

type Operator models.Operator

type Service struct {
	collection *mongo.Collection
}

func NewService(clients *clients.Clients) *Service {
	return &Service{
		collection: clients.Mongo.Db.Collection(constants.CollectionOperators),
	}
}

func (s *Service) AddOperator(ctx context.Context, dto dto.CreateOperatorDTO) (*Operator, error) {
	logger.Info(ctx, "Adding new operator: "+dto.Email)

	t := time.Now().UTC()

	op := &Operator{
		ID:        ulid.MustNew(uint64(t.Unix()), ulid.DefaultEntropy()).String(),
		Username:  dto.Username,
		Email:     dto.Email,
		PwdHash:   dto.PwdHash,
		Role:      dto.Role,
		CreatedAt: t,
	}

	_, err := s.collection.InsertOne(ctx, op)
	if mongo.IsDuplicateKeyError(err) {
		logger.Info(ctx, "Operator already exists: "+dto.Email)
		return nil, ErrOperatorAlreadyExists
	}
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to insert Operator: %w", err))
		return nil, fmt.Errorf("failed to insert Operator: %w", err)
	}

	logger.Info(ctx, "Operator successfully added: "+op.ID)
	return op, nil
}

func (s *Service) GetOperatorByEmail(ctx context.Context, email string) (*Operator, error) {
	logger.Info(ctx, "Fetching operator by email: "+email)

	var op *Operator
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&op)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to find Operator by Email: %w", err))
		return nil, fmt.Errorf("failed to find Operator by Email: %w", err)
	}
	return op, nil
}

func (s *Service) GetOperatorByID(ctx context.Context, id string) (*Operator, error) {
	logger.Info(ctx, "Fetching operator by ID: "+id)

	var op *Operator
	err := s.collection.FindOne(ctx, bson.M{"id": id}).Decode(&op)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to find Operator by ID: %w", err))
		return nil, fmt.Errorf("failed to find Operator by ID: %w", err)
	}
	return op, nil
}

func (s *Service) GetAllOperators(ctx context.Context) ([]*Operator, error) {
	logger.Info(ctx, "Fetching all operators")

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to fetch operators: %w", err))
		return nil, fmt.Errorf("failed to fetch operators: %w", err)
	}
	defer cursor.Close(ctx)

	var result []*Operator
	for cursor.Next(ctx) {
		var op Operator
		if err := cursor.Decode(&op); err != nil {
			logger.Error(ctx, fmt.Errorf("failed to decode operator: %w", err))
			return nil, err
		}
		result = append(result, &op)
	}
	return result, nil
}
