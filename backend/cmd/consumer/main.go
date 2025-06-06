package main

import (
	"context"
	"course_project/cmd/consumer/handlers"
	"course_project/internal/clients"
	"course_project/internal/config"
	"course_project/internal/constants/consumer"
	"course_project/internal/services/logger"
	"encoding/json"
	"fmt"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		logger.Fatal(ctx, fmt.Errorf("config load error: %w", err))
	}

	clnts, err := clients.NewClients(ctx, cfg)
	if err != nil {
		logger.Fatal(ctx, fmt.Errorf("clients init error: %w", err))
	}

	ch := clnts.RabbitMQ.Channel

	msgs, err := ch.Consume(
		cfg.RabbitMQQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Fatal(ctx, fmt.Errorf("failed to register consumer: %w", err))
	}

	logger.Info(ctx, "Message consumer started...")

	handler := handlers.NewHandler(clnts)

	for d := range msgs {
		var env consumer.Envelope
		if err := json.Unmarshal(d.Body, &env); err != nil {
			logger.Error(ctx, fmt.Errorf("failed to parse envelope: %w", err))
			continue
		}

		handler.Handle(ctx, env)
	}
}
