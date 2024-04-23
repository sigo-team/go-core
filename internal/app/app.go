package app

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiber_logger "github.com/gofiber/fiber/v2/middleware/logger"
	"log/slog"
	"os"
	"sigo/internal/config"
	"sigo/internal/controllers"
	"sigo/internal/services"
	http_server "sigo/internal/transport"
)

type App struct {
	app *fiber.App
	cfg *config.Config
}

func New(cfg *config.Config) *App {
	userService := services.NewUserService()

	app := fiber.New(fiber.Config{
		BodyLimit: 256 * 1024 * 1024,
	})

	app.Use(http_server.AuthMiddleware(userService, cfg))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "*",
		AllowCredentials: true,
		//AllowHeaders:     "*",
	}))
	app.Use(fiber_logger.New())

	monoService := services.NewMonoService()

	http_server.PublicRoutes(app, &controllers.RoomHandlers{Service: monoService})

	return &App{
		app: app,
		cfg: cfg,
	}
}

func (a *App) Run() {
	var (
		err     error
		address = fmt.Sprintf("%s:%d", a.cfg.Host, a.cfg.Port)
	)

	log.Info(fmt.Sprintf("Listening on %s", address))

	err = a.app.Listen(address)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func (a *App) Stop() {
	var (
		err error
	)
	log.Info("Gracefully shutting down")
	err = a.app.Shutdown()
	if err != nil {
		slog.Error(err.Error())
	}
	log.Info("Gracefully stopped")
}
