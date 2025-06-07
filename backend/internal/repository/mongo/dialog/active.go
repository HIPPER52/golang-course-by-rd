package dialog

import (
	"context"
	"course_project/internal/models"
	"course_project/internal/repository"
	repo "course_project/internal/repository/dialog/active"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type activeRepo struct {
	col *mongo.Collection
}

func NewActiveRepo(col *mongo.Collection) repo.Repository {
	return &activeRepo{col: col}
}

func (r *activeRepo) Add(ctx context.Context, dialog *models.ActiveDialog) error {
	filter := bson.M{"id": dialog.ID}
	opts := options.Update().SetUpsert(true)
	_, err := r.col.UpdateOne(ctx, filter, bson.M{"$setOnInsert": dialog}, opts)
	if mongo.IsDuplicateKeyError(err) {
		return repository.ErrActiveDialogExists
	}
	return err
}

func (r *activeRepo) ListAll(ctx context.Context) ([]models.ActiveDialog, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dialogs []models.ActiveDialog
	if err := cursor.All(ctx, &dialogs); err != nil {
		return nil, err
	}
	return dialogs, nil
}

func (r *activeRepo) FindByID(ctx context.Context, id string) (*models.ActiveDialog, error) {
	var dialog models.ActiveDialog
	err := r.col.FindOne(ctx, bson.M{"id": id}).Decode(&dialog)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &dialog, nil
}

func (r *activeRepo) FindByOperatorID(ctx context.Context, operatorID string) ([]models.ActiveDialog, error) {
	cursor, err := r.col.Find(ctx, bson.M{"operator_id": operatorID})
	if err != nil {
		return nil, err
	}
	var dialogs []models.ActiveDialog
	if err := cursor.All(ctx, &dialogs); err != nil {
		return nil, err
	}
	return dialogs, nil
}

func (r *activeRepo) CountByOperator(ctx context.Context, operatorID string) (int, error) {
	count, err := r.col.CountDocuments(ctx, bson.M{"operator_id": operatorID})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
