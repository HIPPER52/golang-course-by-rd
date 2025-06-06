package services

import (
	"course_project/internal/clients"
	"course_project/internal/config"
	"course_project/internal/services/auth"
	"course_project/internal/services/client"
	"course_project/internal/services/dialogs/active"
	"course_project/internal/services/dialogs/archived"
	"course_project/internal/services/dialogs/mover"
	"course_project/internal/services/dialogs/queued"
	"course_project/internal/services/message"
	"course_project/internal/services/operator"
	"course_project/internal/services/producer"
)

type Services struct {
	Auth           *auth.Service
	Operator       *operator.Service
	Client         *client.Service
	QueuedDialog   *queued.Service
	ActiveDialog   *active.Service
	ArchivedDialog *archived.Service
	Message        *message.Service
	Mover          *mover.Service
	Producer       *producer.Service
}

func NewServices(cfg *config.Config, clients *clients.Clients) *Services {
	return &Services{
		Auth:           auth.NewService(cfg),
		Operator:       operator.NewService(clients),
		Client:         client.NewService(clients),
		QueuedDialog:   queued.NewService(clients),
		ActiveDialog:   active.NewService(clients),
		ArchivedDialog: archived.NewService(clients),
		Message:        message.NewService(clients),
		Mover:          mover.NewService(clients),
		Producer:       producer.NewService(clients),
	}
}
