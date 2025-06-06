package mover

import (
	"context"
	"course_project/internal/clients"
	"course_project/internal/constants"
	"course_project/internal/models"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	queuedColl   *mongo.Collection
	activeColl   *mongo.Collection
	archivedColl *mongo.Collection
}

func NewService(clients *clients.Clients) *Service {
	return &Service{
		queuedColl:   clients.Mongo.Db.Collection(constants.CollectionQueuedDialog),
		activeColl:   clients.Mongo.Db.Collection(constants.CollectionActiveDialog),
		archivedColl: clients.Mongo.Db.Collection(constants.CollectionArchivedDialog),
	}
}

func (m *Service) TakeDialog(ctx context.Context, dialogID, operatorID string) error {
	var queued models.QueuedDialog
	filter := bson.M{"id": dialogID}
	err := m.queuedColl.FindOne(ctx, filter).Decode(&queued)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("dialog %q not found in queue", dialogID)
		}
		return err
	}

	active := models.ActiveDialog{
		DialogBase: models.DialogBase{
			ID:            queued.ID,
			ClientID:      queued.ClientID,
			ClientName:    queued.ClientName,
			ClientPhone:   queued.ClientPhone,
			ClientIP:      queued.ClientIP,
			OperatorID:    operatorID,
			StartedAt:     time.Now().UTC(),
			LastMessageAt: queued.StartedAt,
		},
	}

	if _, err := m.activeColl.InsertOne(ctx, active); err != nil {
		return fmt.Errorf("failed to insert ActiveDialog: %w", err)
	}

	if _, err := m.queuedColl.DeleteOne(ctx, filter); err != nil {
		return fmt.Errorf("failed to delete from queued-dialog: %w", err)
	}

	return nil
}

func (m *Service) CloseDialog(ctx context.Context, dialogID string) error {
	var active models.ActiveDialog
	filter := bson.M{"id": dialogID}
	err := m.activeColl.FindOne(ctx, filter).Decode(&active)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("dialog %q not found in active", dialogID)
		}
		return err
	}

	archived := models.ArchivedDialog{
		DialogBase: models.DialogBase{
			ID:            active.ID,
			ClientID:      active.ClientID,
			ClientName:    active.ClientName,
			ClientPhone:   active.ClientPhone,
			ClientIP:      active.ClientIP,
			OperatorID:    active.OperatorID,
			StartedAt:     active.StartedAt,
			LastMessageAt: active.LastMessageAt,
			EndedAt:       time.Now().UTC(),
		},
	}

	if _, err := m.archivedColl.InsertOne(ctx, archived); err != nil {
		return fmt.Errorf("failed to insert archived-dialog: %w", err)
	}

	if _, err := m.activeColl.DeleteOne(ctx, filter); err != nil {
		return fmt.Errorf("failed to delete from active-dialog: %w", err)
	}

	return nil
}
