package archived

import (
	"context"
	"course_project/internal/clients"
	"course_project/internal/constants"
	"course_project/internal/models"
	"errors"
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
	filter := bson.M{"id": dialog.ID}
	opts := options.Update().SetUpsert(true)
	_, err := s.collection.UpdateOne(ctx, filter, bson.M{"$setOnInsert": dialog}, opts)
	if mongo.IsDuplicateKeyError(err) {
		return ErrArchivedDialogExists
	}
	return err
}

func (s *Service) FindByOperator(ctx context.Context, operatorID string) ([]models.ArchivedDialog, error) {
	cursor, err := s.collection.Find(ctx, bson.M{"operator_id": operatorID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dialogs []models.ArchivedDialog
	if err := cursor.All(ctx, &dialogs); err != nil {
		return nil, err
	}

	return dialogs, nil
}

func (s *Service) CountByOperator(ctx context.Context, operatorID string) (int, error) {
	count, err := s.collection.CountDocuments(ctx, bson.M{"operator_id": operatorID})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
