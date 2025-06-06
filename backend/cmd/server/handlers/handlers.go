package handlers

import (
	"course_project/cmd/server/handlers/auth"
	"course_project/cmd/server/handlers/client"
	"course_project/cmd/server/handlers/common"
	"course_project/cmd/server/handlers/health"
	"course_project/cmd/server/handlers/operator"
	"course_project/cmd/server/handlers/ws"
	"course_project/cmd/server/middlewares"
	"course_project/internal/config"
	"course_project/internal/constants/roles"
	"course_project/internal/services"
	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	Health   *health.Handler
	Auth     *auth.Handler
	Client   *client.Handler
	Common   *common.Handler
	Operator *operator.Handler
	WS       *ws.Handler

	mdlwrs *middlewares.Middlewares
}

func NewHandlers(
	cfg *config.Config,
	svcs *services.Services,
	mdlwrs *middlewares.Middlewares,
) *Handlers {
	ws := ws.NewHandler(svcs)
	return &Handlers{
		Health:   health.NewHandler(cfg),
		Auth:     auth.NewHandler(svcs),
		Client:   client.NewHandler(svcs, ws),
		Operator: operator.NewHandler(svcs),
		Common:   common.NewHandler(svcs),
		WS:       ws,
		mdlwrs:   mdlwrs,
	}
}

func (h *Handlers) RegisterRoutes(router fiber.Router) {
	//
	//router.Get("/docs/*")
	router.Get("/health", h.Health.Health)

	api := router.Group("/api")
	api.Use(h.mdlwrs.Logger.Handle)

	authGroup := api.Group("/auth")
	authGroup.Post("/signup", h.Auth.SignUp)
	authGroup.Post("/signin", h.Auth.SignIn)

	adminGroup := api.Group("/admin")
	adminGroup.Use(h.mdlwrs.Auth.Handle)
	adminGroup.Use(h.mdlwrs.Role.RequireRoles(roles.Admin))

	operatorGroup := api.Group("/operator")
	operatorGroup.Use(h.mdlwrs.Auth.Handle)
	operatorGroup.Use(h.mdlwrs.Role.RequireRoles(roles.Admin, roles.Operator))
	operatorGroup.Get("/dialogs/queued", h.Operator.GetQueuedDialogs)
	operatorGroup.Get("/dialogs/active", h.Operator.GetActiveDialogs)
	operatorGroup.Get("/dialogs/:room_id/messages", h.Common.GetDialogMessages)

	clientGroup := api.Group("/client")
	clientGroup.Post("/register", h.Client.Register)
	clientGroup.Get("/dialogs/:room_id/messages", h.Common.GetDialogMessages)

	wsGroup := api.Group("/ws")
	wsGroup.Get("/client", h.WS.ClientHandler.Handle())

	wsGroup.Use(h.mdlwrs.Auth.Handle)
	wsGroup.Get("/operator", h.WS.OperatorHandler.Handle())

}
