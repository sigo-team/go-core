package app

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiber_logger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/goombaio/namegenerator"
	"log/slog"
	"os"
	"sigo/internal/config"
	"sigo/internal/controllers"
	"sigo/internal/lib"
	"sigo/internal/services"
	http_server "sigo/internal/transport"
	"time"
)

type App struct {
	app *fiber.App
	cfg *config.Config
	//cancelFunc context.CancelFunc
}

func New(cfg *config.Config) *App {
	app := fiber.New(fiber.Config{
		BodyLimit: 256 * 1024 * 1024,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "*",
		AllowCredentials: true,
	}))
	app.Use(fiber_logger.New())

	roomService := services.NewRoomService(
		services.RoomServiceOptions{
			IdentifierManager: lib.NewIdentifierManager(),
		},
	)
	userService := services.NewUserService(
		services.UserServiceOptions{
			IdentifierManager: lib.NewIdentifierManager(),
			NameGenerator:     namegenerator.NewNameGenerator(time.Now().UnixNano()),
		},
	)

	app.Use(http_server.AuthMiddleware(userService, cfg))

	//closingCtx, closeCtx := context.WithCancel(context.Background())

	roomController := controllers.NewRoomController(
		controllers.RoomControllerOptions{RoomService: roomService},
	)

	http_server.PublicRoutes(app, roomController)

	return &App{
		app: app,
		cfg: cfg,
		//cancelFunc: closeCtx,
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
