package archived

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
	ErrArchivedDialogExists = errors.New("archived dialog already exists")
)

type Service struct {
	collection *mongo.Collection
}

func NewService(clients *clients.Clients) *Service {
	return &Service{collection: clients.Mongo.Db.Collection(constants.CollectionArchivedDialog)}
}

func (s *Service) Add(ctx context.Context, dialog *models.ArchivedDialog) error {
	logger.Info(ctx, "Archiving dialog: "+dialog.ID)

	filter := bson.M{"id": dialog.ID}
	opts := options.Update().SetUpsert(true)
	_, err := s.collection.UpdateOne(ctx, filter, bson.M{"$setOnInsert": dialog}, opts)
	if mongo.IsDuplicateKeyError(err) {
		logger.Error(ctx, fmt.Errorf("failed to archive dialog: %w", err))
		return ErrArchivedDialogExists
	}
	return err
}

func (s *Service) FindByOperator(ctx context.Context, operatorID string) ([]models.ArchivedDialog, error) {
	logger.Info(ctx, "Finding archived dialogs by operator: "+operatorID)

	cursor, err := s.collection.Find(ctx, bson.M{"operator_id": operatorID})
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to find archived dialogs: %w", err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var dialogs []models.ArchivedDialog
	if err := cursor.All(ctx, &dialogs); err != nil {
		logger.Error(ctx, fmt.Errorf("failed to decode archived dialogs: %w", err))
		return nil, err
	}

	return dialogs, nil
}

func (s *Service) CountByOperator(ctx context.Context, operatorID string) (int, error) {
	logger.Info(ctx, "Counting archived dialogs by operator: "+operatorID)

	count, err := s.collection.CountDocuments(ctx, bson.M{"operator_id": operatorID})
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to count archived dialogs: %w", err))
		return 0, err
	}
	return int(count), nil
}
