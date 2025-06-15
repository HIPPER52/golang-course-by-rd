package dialog

import (
	"context"
	"course_project/internal/models"
	"course_project/internal/repository"
	"course_project/internal/repository/dialog/queued"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ queued.Repository = (*queuedRepo)(nil)

type queuedRepo struct {
	col *mongo.Collection
}

func NewQueuedRepo(col *mongo.Collection) queued.Repository {
	return &queuedRepo{col: col}
}

func (r *queuedRepo) Add(ctx context.Context, dialog *models.QueuedDialog) error {
	filter := bson.M{"id": dialog.ID}
	opts := options.Update().SetUpsert(true)
	_, err := r.col.UpdateOne(ctx, filter, bson.M{"$setOnInsert": dialog}, opts)
	if mongo.IsDuplicateKeyError(err) {
		return repository.ErrQueuedDialogExists
	}
	return err
}

func (r *queuedRepo) FindByID(ctx context.Context, id string) (*models.QueuedDialog, error) {
	var dialog models.QueuedDialog
	err := r.col.FindOne(ctx, bson.M{"id": id}).Decode(&dialog)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &dialog, nil
}

func (r *queuedRepo) ListAll(ctx context.Context) ([]models.QueuedDialog, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dialogs []models.QueuedDialog
	if err := cursor.All(ctx, &dialogs); err != nil {
		return nil, err
	}
	return dialogs, nil
}
