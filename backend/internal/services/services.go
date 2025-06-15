package services

import (
	"course_project/internal/clients"
	"course_project/internal/config"
	"course_project/internal/constants"
	mongoRepo "course_project/internal/repository/mongo"
	dialogRepo "course_project/internal/repository/mongo/dialog"
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
	db := clients.Mongo.Db
	clientRepo := mongoRepo.NewClientRepo(db.Collection(constants.CollectionClients))
	messageRepo := mongoRepo.NewMessageRepo(db.Collection(constants.CollectionMessages))
	operatorRepo := mongoRepo.NewOperatorRepo(db.Collection(constants.CollectionOperators))
	dialogFinder := mongoRepo.NewDialogFinder(
		db.Collection(constants.CollectionQueuedDialog),
		db.Collection(constants.CollectionActiveDialog),
		db.Collection(constants.CollectionArchivedDialog),
	)
	moverRepo := mongoRepo.NewDialogMoverRepo(
		db.Collection(constants.CollectionQueuedDialog),
		db.Collection(constants.CollectionActiveDialog),
		db.Collection(constants.CollectionArchivedDialog),
	)
	activeRepo := dialogRepo.NewActiveRepo(db.Collection(constants.CollectionActiveDialog))
	archivedRepo := dialogRepo.NewArchivedRepo(db.Collection(constants.CollectionArchivedDialog))
	queuedRepo := dialogRepo.NewQueuedRepo(db.Collection(constants.CollectionQueuedDialog))

	return &Services{
		Auth:           auth.NewService(cfg),
		Operator:       operator.NewService(operatorRepo),
		Client:         client.NewService(clientRepo),
		QueuedDialog:   queued.NewService(queuedRepo),
		ActiveDialog:   active.NewService(activeRepo),
		ArchivedDialog: archived.NewService(archivedRepo),
		Message:        message.NewService(messageRepo, dialogFinder),
		Mover:          mover.NewService(moverRepo),
		Producer:       producer.NewService(clients),
	}
}
