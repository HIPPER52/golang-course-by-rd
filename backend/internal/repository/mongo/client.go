package mongo

import (
	"context"
	"course_project/internal/models"
	"course_project/internal/repository/client"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type clientRepo struct {
	col *mongo.Collection
}

func NewClientRepo(col *mongo.Collection) client.Repository {
	return &clientRepo{col: col}
}

func (r *clientRepo) CountByPhone(ctx context.Context, phone string) (int64, error) {
	return r.col.CountDocuments(ctx, bson.M{"phone": phone})
}

func (r *clientRepo) Create(ctx context.Context, c *models.Client) error {
	_, err := r.col.InsertOne(ctx, c)
	return err
}
