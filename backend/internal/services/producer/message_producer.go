package producer

import (
	"course_project/internal/clients"
	"course_project/internal/clients/rabbitmq"
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
	body, err := json.Marshal(Envelope{
		Type:    eventType,
		Payload: payload,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

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
