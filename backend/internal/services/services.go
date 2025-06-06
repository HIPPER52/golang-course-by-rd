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
	"course_project/internal/services/operator"
)

type Services struct {
	Auth           *auth.Service
	Operator       *operator.Service
	Client         *client.Service
	QueuedDialog   *queued.Service
	ActiveDialog   *active.Service
	ArchivedDialog *archived.Service
	Mover          *mover.Service
}

func NewServices(cfg *config.Config, clients *clients.Clients) *Services {
	return &Services{
		Auth:           auth.NewService(cfg),
		Operator:       operator.NewService(clients),
		Client:         client.NewService(clients),
		QueuedDialog:   queued.NewService(clients),
		ActiveDialog:   active.NewService(clients),
		ArchivedDialog: archived.NewService(clients),
		Mover:          mover.NewService(clients),
	}
}
