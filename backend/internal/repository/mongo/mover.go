package mongo

import (
	"context"
	"course_project/internal/models"
	"course_project/internal/repository/dialog/mover"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type dialogMoverRepo struct {
	queued   *mongo.Collection
	active   *mongo.Collection
	archived *mongo.Collection
}

func NewDialogMoverRepo(queued, active, archived *mongo.Collection) mover.Repository {
	return &dialogMoverRepo{queued, active, archived}
}

func (r *dialogMoverRepo) FindQueuedByID(ctx context.Context, id string) (*models.QueuedDialog, error) {
	var q models.QueuedDialog
	err := r.queued.FindOne(ctx, bson.M{"id": id}).Decode(&q)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &q, nil
}

func (r *dialogMoverRepo) InsertActive(ctx context.Context, dialog models.ActiveDialog) error {
	_, err := r.active.InsertOne(ctx, dialog)
	return err
}

func (r *dialogMoverRepo) DeleteQueuedByID(ctx context.Context, id string) error {
	_, err := r.queued.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (r *dialogMoverRepo) FindActiveByID(ctx context.Context, id string) (*models.ActiveDialog, error) {
	var a models.ActiveDialog
	err := r.active.FindOne(ctx, bson.M{"id": id}).Decode(&a)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *dialogMoverRepo) InsertArchived(ctx context.Context, dialog models.ArchivedDialog) error {
	_, err := r.archived.InsertOne(ctx, dialog)
	return err
}

func (r *dialogMoverRepo) DeleteActiveByID(ctx context.Context, id string) error {
	_, err := r.active.DeleteOne(ctx, bson.M{"id": id})
	return err
}
