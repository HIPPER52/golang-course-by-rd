package message

import (
	"context"
	"course_project/internal/clients"
	"course_project/internal/constants"
	"course_project/internal/models"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrAccessDenied   = errors.New("you don't have access to this dialog")
	ErrDialogNotFound = errors.New("dialog not found")
)

type Service struct {
	messagesColl        *mongo.Collection
	archivedDialogsColl *mongo.Collection
	activeDialogsColl   *mongo.Collection
	queuedDialogsColl   *mongo.Collection
}

func NewService(clients *clients.Clients) *Service {
	return &Service{
		messagesColl:        clients.Mongo.Db.Collection(constants.CollectionMessages),
		archivedDialogsColl: clients.Mongo.Db.Collection(constants.CollectionArchivedDialog),
		activeDialogsColl:   clients.Mongo.Db.Collection(constants.CollectionActiveDialog),
		queuedDialogsColl:   clients.Mongo.Db.Collection(constants.CollectionQueuedDialog),
	}
}

func (s *Service) FindByRoomID(ctx context.Context, roomID string, clientID string) ([]models.Message, error) {
	dialog, err := s.findDialogByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	if clientID != "" && dialog.ClientID != clientID {
		return nil, ErrAccessDenied
	}

	opts := options.Find().SetSort(bson.M{"sent_at": 1})
	cursor, err := s.messagesColl.Find(ctx, bson.M{"room_id": roomID}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find messages: %w", err)
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, fmt.Errorf("failed to decode messages: %w", err)
	}

	return messages, nil
}

func (s *Service) findDialogByID(ctx context.Context, roomID string) (*models.ArchivedDialog, error) {
	colls := []*mongo.Collection{
		s.queuedDialogsColl,
		s.activeDialogsColl,
		s.archivedDialogsColl,
	}

	for _, coll := range colls {
		var dialog models.ArchivedDialog
		err := coll.FindOne(ctx, bson.M{"id": roomID}).Decode(&dialog)
		if err == nil {
			return &dialog, nil
		}
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("failed to fetch dialog: %w", err)
		}
	}

	return nil, ErrDialogNotFound
}
