package operator

import (
	"context"
	"course_project/internal/clients"
	"course_project/internal/constants"
	"course_project/internal/dto"
	"course_project/internal/models"
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
		return nil, ErrOperatorAlreadyExists
	}
	if err != nil {
		return nil, fmt.Errorf("failed to insert Operator: %w", err)
	}

	return op, nil
}

func (s *Service) GetOperatorByEmail(ctx context.Context, email string) (*Operator, error) {
	var op *Operator
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&op)
	if err != nil {
		return nil, fmt.Errorf("failed to find Operator by Email: %w", err)
	}
	return op, nil
}

func (s *Service) GetOperatorByID(ctx context.Context, id string) (*Operator, error) {
	var op *Operator
	err := s.collection.FindOne(ctx, bson.M{"id": id}).Decode(&op)
	if err != nil {
		return nil, fmt.Errorf("failed to find Operator by ID: %w", err)
	}
	return op, nil
}
