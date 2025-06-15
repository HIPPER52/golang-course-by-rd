// @title           Support Chat Server
// @version         1.0
// @description     Course project: simple support chat with roles, dialogs, and messaging.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Dima Avtenev
// @contact.email  hipper52@gmail.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey X-User-Token
// @in header
// @name Authorization
package main

import (
	"context"
	_ "course_project/cmd/server/docs"
	"course_project/cmd/server/handlers"
	"course_project/cmd/server/middlewares"
	"course_project/internal/clients"
	"course_project/internal/config"
	"course_project/internal/services"
	"course_project/internal/services/logger"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

func main() {
	ctx := context.Background()
	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		logger.Fatal(ctx, fmt.Errorf("failed to load config: %v", err))
	}

	clnts, err := clients.NewClients(ctx, cfg)
	if err != nil {
		logger.Fatal(ctx, fmt.Errorf("failed to create clients: %v", err))
	}

	svcs := services.NewServices(cfg, clnts)
	mdlwrs := middlewares.NewMiddlewares(svcs)

	app := fiber.New()

	app.Use(recover.New())

	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, X-User-Token",
	}))

	hdlrs := handlers.NewHandlers(cfg, svcs, mdlwrs)
	hdlrs.RegisterRoutes(app)

	app.Get("/swagger/*", swagger.HandlerDefault)

	go func() {
		err := app.Listen(":" + cfg.Port)
		if err != nil {
			logger.Panic(ctx, fmt.Errorf("server listening failed: %v", err))
		}
	}()

	logger.Info(ctx, "Server started...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	err = app.Shutdown()
	if err != nil {
		logger.Fatal(ctx, fmt.Errorf("server shutdown failed: %v", err))
	}
}
