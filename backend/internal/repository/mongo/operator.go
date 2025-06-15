package mongo

import (
	"context"
	"course_project/internal/dto"
	"course_project/internal/models"
	"course_project/internal/repository"
	"course_project/internal/repository/operator"
	"fmt"
	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type operatorRepo struct {
	col *mongo.Collection
}

func NewOperatorRepo(col *mongo.Collection) operator.Repository {
	return &operatorRepo{col: col}
}

func (r *operatorRepo) AddOperator(ctx context.Context, dto dto.CreateOperatorDTO) (*models.Operator, error) {
	t := time.Now().UTC()

	op := &models.Operator{
		ID:        ulid.MustNew(uint64(t.Unix()), ulid.DefaultEntropy()).String(),
		Username:  dto.Username,
		Email:     dto.Email,
		PwdHash:   dto.PwdHash,
		Role:      dto.Role,
		CreatedAt: t,
	}

	_, err := r.col.InsertOne(ctx, op)
	if mongo.IsDuplicateKeyError(err) {
		return nil, repository.ErrOperatorAlreadyExists
	}
	if err != nil {
		return nil, fmt.Errorf("failed to insert Operator: %w", err)
	}

	return op, nil
}

func (r *operatorRepo) GetOperatorByEmail(ctx context.Context, email string) (*models.Operator, error) {
	var op *models.Operator
	err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&op)
	if err != nil {
		return nil, fmt.Errorf("failed to find Operator by Email: %w", err)
	}
	return op, nil
}

func (r *operatorRepo) GetOperatorByID(ctx context.Context, id string) (*models.Operator, error) {
	var op *models.Operator
	err := r.col.FindOne(ctx, bson.M{"id": id}).Decode(&op)
	if err != nil {
		return nil, fmt.Errorf("failed to find Operator by ID: %w", err)
	}
	return op, nil
}

func (r *operatorRepo) GetAllOperators(ctx context.Context) ([]*models.Operator, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch operators: %w", err)
	}
	defer cursor.Close(ctx)

	var result []*models.Operator
	for cursor.Next(ctx) {
		var op models.Operator
		if err := cursor.Decode(&op); err != nil {
			return nil, err
		}
		result = append(result, &op)
	}
	return result, nil
}
