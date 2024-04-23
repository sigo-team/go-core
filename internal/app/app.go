package app

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiber_logger "github.com/gofiber/fiber/v2/middleware/logger"
	"go-core/internal/config"
	"log/slog"
	"os"
)

type App struct {
	app *fiber.App
}

func New() *App {
	var (
		service *mono.MonoService
		app     *fiber.App
		//db      *storage.MainDB
	)

	//db = storage.New("app.db")

	service = mono.New()
	authService := auth.New()
	userService := user.New()

	app = fiber.New(fiber.Config{
		BodyLimit: 256 * 1024 * 1024,
	})

	app.Use(middleware.AuthMiddleware(authService, userService))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "*",
		AllowCredentials: true,
		//AllowHeaders:     "*",
	}))
	app.Use(fiber_logger.New())

	routes.PublicRoutes(app, &handlers.RoomHandlers{
		Service: service,
	})

	return &App{
		app: app,
		//db:  db,
	}
}

func (a *App) Run() {
	var (
		cfg     = config.Cfg()
		err     error
		address = fmt.Sprintf("%s:%d", cfg.HOST, cfg.PORT)
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
	//err = a.db.Close()
	//if err != nil {
	//  log.Error(err.Error())
	//}
	log.Info("Gracefully shutting down")
	err = a.app.Shutdown()
	if err != nil {
		slog.Error(err.Error())
	}
	log.Info("Gracefully stopped")
}
