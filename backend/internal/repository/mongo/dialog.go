package mongo

import (
	"context"
	"course_project/internal/models"
	repo "course_project/internal/repository/message"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type dialogFinder struct {
	queued   *mongo.Collection
	active   *mongo.Collection
	archived *mongo.Collection
}

var _ repo.DialogFinder = (*dialogFinder)(nil)

func NewDialogFinder(queued, active, archived *mongo.Collection) *dialogFinder {
	return &dialogFinder{
		queued:   queued,
		active:   active,
		archived: archived,
	}
}

func (d *dialogFinder) FindDialogByID(ctx context.Context, roomID string) (*models.ArchivedDialog, error) {
	colls := []*mongo.Collection{d.queued, d.active, d.archived}

	for _, coll := range colls {
		var dialog models.ArchivedDialog
		err := coll.FindOne(ctx, bson.M{"id": roomID}).Decode(&dialog)
		if err == nil {
			return &dialog, nil
		}
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("failed to fetch dialog from %s: %w", coll.Name(), err)
		}
	}

	return nil, errors.New("dialog not found")
}
