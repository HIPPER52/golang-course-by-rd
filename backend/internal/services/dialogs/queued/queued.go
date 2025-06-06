package queued

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
	ErrQueuedDialogExists = errors.New("queued dialog already exists")
)

type Service struct {
	collection *mongo.Collection
}

func NewService(clients *clients.Clients) *Service {
	return &Service{collection: clients.Mongo.Db.Collection(constants.CollectionQueuedDialog)}
}

func (s *Service) Add(ctx context.Context, dialog *models.QueuedDialog) error {
	filter := bson.M{"id": dialog.ID}
	opts := options.Update().SetUpsert(true)
	_, err := s.collection.UpdateOne(ctx, filter, bson.M{"$setOnInsert": dialog}, opts)
	if mongo.IsDuplicateKeyError(err) {
		return ErrQueuedDialogExists
	}
	return err
}

func (s *Service) FindByID(ctx context.Context, id string) (*models.QueuedDialog, error) {
	var dialog models.QueuedDialog
	err := s.collection.FindOne(ctx, map[string]any{"id": id}).Decode(&dialog)
	if err != nil {
		return nil, err
	}
	return &dialog, nil
}

func (s *Service) ListAll(ctx context.Context) ([]models.QueuedDialog, error) {
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dialogs []models.QueuedDialog
	for cursor.Next(ctx) {
		var d models.QueuedDialog
		if err := cursor.Decode(&d); err != nil {
			return nil, err
		}
		dialogs = append(dialogs, d)
	}

	return dialogs, nil
}
