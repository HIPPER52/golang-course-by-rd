package producer

import (
	"course_project/internal/clients"
	"course_project/internal/clients/rabbitmq"
	"course_project/internal/services/logger"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
)

type Envelope struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type Service struct {
	rmq      *rabbitmq.Client
	exchange string
}

func NewService(clients *clients.Clients) *Service {
	return &Service{
		rmq:      clients.RabbitMQ,
		exchange: clients.RabbitMQ.Exchange,
	}
}

func (p *Service) Publish(eventType string, payload interface{}) error {
	logger.Info(nil, "Publishing event: "+eventType)

	body, err := json.Marshal(Envelope{
		Type:    eventType,
		Payload: payload,
	})
	if err != nil {
		logger.Error(nil, fmt.Errorf("failed to marshal event %s: %w", eventType, err))
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	logger.Info(nil, "Event published: "+eventType)
	return p.rmq.Channel.Publish(
		p.exchange,
		"",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
