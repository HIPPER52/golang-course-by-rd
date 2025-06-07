package dialog

import (
	"context"
	"course_project/internal/models"
	"course_project/internal/repository/dialog/archived"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ archived.Repository = (*archivedRepo)(nil)

var ErrArchivedDialogExists = errors.New("archived dialog already exists")

type archivedRepo struct {
	col *mongo.Collection
}

func NewArchivedRepo(col *mongo.Collection) *archivedRepo {
	return &archivedRepo{col: col}
}

func (r *archivedRepo) Add(ctx context.Context, dialog *models.ArchivedDialog) error {
	filter := bson.M{"id": dialog.ID}
	opts := options.Update().SetUpsert(true)
	_, err := r.col.UpdateOne(ctx, filter, bson.M{"$setOnInsert": dialog}, opts)
	if mongo.IsDuplicateKeyError(err) {
		return ErrArchivedDialogExists
	}
	return err
}

func (r *archivedRepo) FindByOperator(ctx context.Context, operatorID string) ([]models.ArchivedDialog, error) {
	cursor, err := r.col.Find(ctx, bson.M{"operator_id": operatorID})
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

func (r *archivedRepo) CountByOperator(ctx context.Context, operatorID string) (int, error) {
	count, err := r.col.CountDocuments(ctx, bson.M{"operator_id": operatorID})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
