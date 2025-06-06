package active

import (
	"context"
	"course_project/internal/clients"
	"course_project/internal/constants"
	"course_project/internal/models"
	"course_project/internal/services/logger"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrActiveDialogExists = errors.New("active dialog already exists")
)

type Service struct {
	collection *mongo.Collection
}

func NewService(clients *clients.Clients) *Service {
	return &Service{collection: clients.Mongo.Db.Collection(constants.CollectionActiveDialog)}
}

func (s *Service) Add(ctx context.Context, dialog *models.ActiveDialog) error {
	logger.Info(ctx, "Adding active dialog: "+dialog.ID)

	filter := bson.M{"id": dialog.ID}
	opts := options.Update().SetUpsert(true)
	_, err := s.collection.UpdateOne(ctx, filter, bson.M{"$setOnInsert": dialog}, opts)
	if mongo.IsDuplicateKeyError(err) {
		logger.Info(ctx, "Active dialog already exists: "+dialog.ID)
		return ErrActiveDialogExists
	}
	return err
}

func (s *Service) ListAll(ctx context.Context) ([]models.ActiveDialog, error) {
	logger.Info(ctx, "Listing all active dialogs")

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to list active dialogs: %w", err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var dialogs []models.ActiveDialog
	for cursor.Next(ctx) {
		var d models.ActiveDialog
		if err := cursor.Decode(&d); err != nil {
			logger.Error(ctx, fmt.Errorf("failed to decode active dialog: %w", err))
			return nil, err
		}
		dialogs = append(dialogs, d)
	}

	return dialogs, nil
}

func (s *Service) FindByID(ctx context.Context, id string) (*models.ActiveDialog, error) {
	logger.Info(ctx, "Finding active dialog by ID: "+id)

	var dialog models.ActiveDialog
	err := s.collection.FindOne(ctx, bson.M{"id": id}).Decode(&dialog)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to find active dialog by ID: %w", err))

		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &dialog, nil
}

func (s *Service) FindByOperatorID(ctx context.Context, operatorID string) ([]models.ActiveDialog, error) {
	logger.Info(ctx, "Finding active dialogs by operator: "+operatorID)

	cursor, err := s.collection.Find(ctx, bson.M{"operator_id": operatorID})
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to find active dialogs: %w", err))
		return nil, err
	}
	var dialogs []models.ActiveDialog
	if err := cursor.All(ctx, &dialogs); err != nil {
		logger.Error(ctx, fmt.Errorf("failed to decode active dialogs: %w", err))
		return nil, err
	}
	return dialogs, nil
}

func (s *Service) CountByOperator(ctx context.Context, operatorID string) (int, error) {
	logger.Info(ctx, "Counting active dialogs by operator: "+operatorID)

	count, err := s.collection.CountDocuments(ctx, bson.M{"operator_id": operatorID})
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to count active dialogs: %w", err))
		return 0, err
	}
	return int(count), nil
}
