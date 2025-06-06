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
