package clients

import (
	"context"
	app_mongo "course_project/internal/clients/mongo"
	"course_project/internal/config"
)

type Clients struct {
	Mongo *app_mongo.Client
}

func NewClients(ctx context.Context, cfg *config.Config) (*Clients, error) {
	mongo, err := app_mongo.NewClient(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &Clients{
		Mongo: mongo,
	}, nil
}
