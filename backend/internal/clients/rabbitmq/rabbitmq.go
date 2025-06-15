package rabbitmq

import (
	"course_project/internal/config"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
)

type Client struct {
	Conn     *amqp091.Connection
	Channel  *amqp091.Channel
	Exchange string
}

func NewClient(cfg *config.Config) (*Client, error) {
	var exchange = cfg.RabbitMQExchange
	var queue = cfg.RabbitMQQueue

	conn, err := amqp091.Dial(cfg.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}

	if err := ch.ExchangeDeclare(
		exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	_, err = ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	err = ch.QueueBind(
		queue,
		"",
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	return &Client{
		Conn:     conn,
		Channel:  ch,
		Exchange: exchange,
	}, nil
}

func (c *Client) Close() error {
	if err := c.Channel.Close(); err != nil {
		return fmt.Errorf("failed to close RabbitMQ channel: %w", err)
	}
	if err := c.Conn.Close(); err != nil {
		return fmt.Errorf("failed to close RabbitMQ connection: %w", err)
	}
	return nil
}
