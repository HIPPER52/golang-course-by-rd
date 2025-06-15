package mongo

import (
	"context"
	"course_project/internal/models"
	repo "course_project/internal/repository/message"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type messageRepo struct {
	col *mongo.Collection
}

var _ repo.Repository = (*messageRepo)(nil)

func NewMessageRepo(col *mongo.Collection) *messageRepo {
	return &messageRepo{col: col}
}

func (r *messageRepo) FindByRoomID(ctx context.Context, roomID string) ([]models.Message, error) {
	opts := options.Find().SetSort(bson.M{"sent_at": 1})
	cursor, err := r.col.Find(ctx, bson.M{"room_id": roomID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var msgs []models.Message
	if err := cursor.All(ctx, &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}
