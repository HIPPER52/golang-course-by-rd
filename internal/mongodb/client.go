package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

var client *mongo.Client

func Init() error {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return fmt.Errorf("MONGO_URI not set in environment")
	}

	ctx := context.Background()
	opts := options.Client()
	opts.ApplyURI(mongoURI)

	var err error
	client, err = mongo.Connect(ctx, opts)
	if err != nil {
		return fmt.Errorf("mongo.Connect failed: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	fmt.Println("MongoDB connected successfully")
	return nil
}

func GetClient() *mongo.Client {
	return client
}
