package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/arif-x/sqlx-mysql-boilerplate/app/http/middleware"
	"github.com/arif-x/sqlx-mysql-boilerplate/config"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/database"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/logger"
	route "github.com/arif-x/sqlx-mysql-boilerplate/route/api"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

// Serve ..
func Serve() {
	appCfg := config.AppCfg()

	logger.SetUpLogger()
	logr := logger.GetLogger()

	// connect to DB
	if err := database.ConnectDB(); err != nil {
		logr.Panicf("failed database setup. error: %v", err)
	}

	// Define Fiber config & app.
	fiberCfg := config.FiberConfig()
	app := fiber.New(fiberCfg)

	// Attach Middlewares.
	middleware.FiberMiddleware(app)

	// Routes.
	route.Auth(app)
	route.Dashboard(app)
	route.Public(app)
	route.FileRoutes(app)
	app.Get("/swagger/*", swagger.HandlerDefault)

	// signal channel to capture system calls
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// start shutdown goroutine
	go func() {
		// capture sigterm and other system call here
		<-sigCh
		logr.Infoln("Shutting down server...")
		_ = app.Shutdown()
	}()

	// start http server
	serverAddr := fmt.Sprintf("%s:%d", appCfg.Host, appCfg.Port)
	if err := app.Listen(serverAddr); err != nil {
		logr.Errorf("Oops... server is not running! error: %v", err)
	}

}
