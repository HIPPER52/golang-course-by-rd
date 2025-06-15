package consumer

import (
	"context"
	"course_project/internal/constants/consumer"
	"course_project/internal/services/logger"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	ch      *amqp.Channel
	queue   string
	handler Handler
}

type Handler interface {
	Handle(ctx context.Context, env consumer.Envelope)
}

func New(ch *amqp.Channel, queue string, handler Handler) *Consumer {
	return &Consumer{ch: ch, queue: queue, handler: handler}
}

func (c *Consumer) Start(ctx context.Context) error {
	msgs, err := c.ch.Consume(
		c.queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		logger.Info(ctx, "Consumer listening...")

		for d := range msgs {
			var env consumer.Envelope
			if err := json.Unmarshal(d.Body, &env); err != nil {
				logger.Error(ctx, fmt.Errorf("failed to parse envelope: %w", err))
				continue
			}
			c.handler.Handle(ctx, env)
		}
	}()

	return nil
}
