package transport

import (
	"github.com/gofiber/fiber/v2"
	"sigo/internal/controllers"
	"sigo/internal/ws"
)

func PublicRoutes(app *fiber.App, r *controllers.RoomHandlers) {
	route := app.Group("/api/v1")

	route.Post("/room", r.CreateRoom)
	route.Get("/room", r.GetRooms)

	app.Get("/", controllers.UpgradeHandler())
	app.Get("/", ws.ConnectPlayerHandler(ctx, lb))
	app.Get("/khil", controllers.UpgradeHandler())
	app.Get("/khil", ws.ConnectKhil(ctx, lb))
}
