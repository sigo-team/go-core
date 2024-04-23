package http_server

import (
	"github.com/gofiber/fiber/v2"
	"sigo/internal/controllers"
	"sigo/internal/services"
	"sigo/internal/ws"
)

func PublicRouters(app *fiber.App) {
	route := app.Group("/api/v1")

	monoService := services.NewMono()
	r := controllers.RoomHandlers{Service: monoService}

	route.Post("/room", r.CreateRoom)
	route.Get("/room", r.GetRooms)

	app.Get("/", controllers.UpgradeHandler())
	app.Get("/", ws.ConnectPlayerHandler(ctx, lb))
	app.Get("/khil", controllers.UpgradeHandler())
	app.Get("/khil", ws.ConnectKhil(ctx, lb))
}
