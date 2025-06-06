package client

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

var ErrClientAlreadyExists = errors.New("client already exists")

type Service struct {
	collection *mongo.Collection
}

func NewService(clients *clients.Clients) *Service {
	return &Service{
		collection: clients.Mongo.Db.Collection(constants.CollectionClients),
	}
}

func (s *Service) RegisterClient(ctx context.Context, dto dto.RegisterClientDTO) (*models.Client, error) {
	t := time.Now().UTC()

	count, err := s.collection.CountDocuments(ctx, bson.M{"phone": dto.Phone})
	if err != nil {
		return nil, fmt.Errorf("failed to check existing clients: %w", err)
	}
	if count > 0 {
		return nil, ErrClientAlreadyExists
	}

	client := &models.Client{
		ID:        ulid.Make().String(),
		Name:      dto.Name,
		Phone:     dto.Phone,
		CreatedAt: t,
	}

	if _, err := s.collection.InsertOne(ctx, client); err != nil {
		return nil, fmt.Errorf("failed to insert client: %w", err)
	}

	return client, nil
}
