package clients

import (
	"context"
	app_mongo "course_project/internal/clients/mongo"
	"course_project/internal/clients/rabbitmq"
	"course_project/internal/config"
)

type Clients struct {
	Mongo    *app_mongo.Client
	RabbitMQ *rabbitmq.Client
}

func NewClients(ctx context.Context, cfg *config.Config) (*Clients, error) {
	mongo, err := app_mongo.NewClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	rmq, err := rabbitmq.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Clients{
		Mongo:    mongo,
		RabbitMQ: rmq,
	}, nil
}
