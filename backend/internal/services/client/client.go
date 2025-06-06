package client

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

	logger.Info(ctx, "Registering new client: "+dto.Phone)

	count, err := s.collection.CountDocuments(ctx, bson.M{"phone": dto.Phone})
	if err != nil {
		logger.Error(nil, fmt.Errorf("failed to check existing clients"))
		return nil, fmt.Errorf("failed to check existing clients: %w", err)
	}
	if count > 0 {
		logger.Info(ctx, "Client already exists: "+dto.Phone)
		return nil, ErrClientAlreadyExists
	}

	client := &models.Client{
		ID:        ulid.Make().String(),
		Name:      dto.Name,
		Phone:     dto.Phone,
		CreatedAt: t,
	}

	if _, err := s.collection.InsertOne(ctx, client); err != nil {
		logger.Error(nil, fmt.Errorf("failed to insert client"))
		return nil, fmt.Errorf("failed to insert client: %w", err)
	}

	logger.Info(ctx, "New client registered: "+client.ID)
	return client, nil
}
