package main

import (
	"context"
	"course_project/cmd/consumer/handlers"
	"course_project/internal/clients"
	"course_project/internal/config"
	"course_project/internal/consumer"
	"course_project/internal/services/logger"
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

	handler := handlers.NewHandler(clnts)

	cons := consumer.New(clnts.RabbitMQ.Channel, cfg.RabbitMQQueue, handler)

	if err := cons.Start(ctx); err != nil {
		logger.Fatal(ctx, fmt.Errorf("consumer start error: %w", err))
	}

	select {}
}
