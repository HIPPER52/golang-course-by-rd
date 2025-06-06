package mover

import (
	"context"
	"course_project/internal/clients"
	"course_project/internal/constants"
	"course_project/internal/models"
	"course_project/internal/services/logger"
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
	logger.Info(ctx, "Operator "+operatorID+" is taking dialog: "+dialogID)

	var queued models.QueuedDialog
	filter := bson.M{"id": dialogID}
	err := m.queuedColl.FindOne(ctx, filter).Decode(&queued)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error(ctx, fmt.Errorf("dialog %s not found in queue", dialogID))
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
		logger.Error(ctx, fmt.Errorf("failed to insert active dialog: %w", err))
		return fmt.Errorf("failed to insert ActiveDialog: %w", err)
	}

	if _, err := m.queuedColl.DeleteOne(ctx, filter); err != nil {
		logger.Error(ctx, fmt.Errorf("failed to delete from queued-dialog: %w", err))
		return fmt.Errorf("failed to delete from queued-dialog: %w", err)
	}

	logger.Info(ctx, "Dialog "+dialogID+" successfully taken by operator "+operatorID)
	return nil
}

func (m *Service) CloseDialog(ctx context.Context, dialogID string) error {
	logger.Info(ctx, "Closing dialog: "+dialogID)

	var active models.ActiveDialog
	filter := bson.M{"id": dialogID}
	err := m.activeColl.FindOne(ctx, filter).Decode(&active)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error(ctx, fmt.Errorf("dialog %s not found in active", dialogID))
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
		logger.Error(ctx, fmt.Errorf("failed to insert archived-dialog: %w", err))
		return fmt.Errorf("failed to insert archived-dialog: %w", err)
	}

	if _, err := m.activeColl.DeleteOne(ctx, filter); err != nil {
		logger.Error(ctx, fmt.Errorf("failed to delete from active-dialog: %w", err))
		return fmt.Errorf("failed to delete from active-dialog: %w", err)
	}

	logger.Info(ctx, "Dialog "+dialogID+" successfully archived and closed")
	return nil
}
