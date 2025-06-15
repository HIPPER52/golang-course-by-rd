package app_mongo

import (
	"context"
	"course_project/internal/config"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	Db *mongo.Database
}

func NewClient(ctx context.Context, cfg *config.Config) (*Client, error) {
	opts := options.Client()
	opts.ApplyURI(cfg.MongoURI)
	conn, err := mongo.Connect(ctx, opts)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	return &Client{
		Db: conn.Database(cfg.MongoDbName),
	}, nil
}
