package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"os"
)

func NewMongoClient(ctx context.Context) (*mongo.Client, error) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI not set in environment")
	}

	opts := options.Client()
	opts.ApplyURI(mongoURI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("mongo.Connect failed: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	slog.Info("MongoDB connected successfully")
	return client, nil
}
