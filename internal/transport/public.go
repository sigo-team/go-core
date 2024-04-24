package transport

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"sigo/internal/controllers"
)

func PublicRoutes(app *fiber.App, ctx context.Context, r *controllers.RoomHandlers) {
	route := app.Group("/api/v1")

	route.Post("/room", r.CreateRoom)

	route.Post("/room", UpgraderMiddleware(), ConnectKhil(ctx, r))

	//todo" query params
	route.Post("/connect")

	route.Get("/room", r.GetRooms)

	//app.Get("/", controllers.UpgradeHandler())
	//app.Get("/", ws.ConnectPlayerHandler(ctx, lb))
	//app.Get("/khil", controllers.UpgradeHandler())
	//app.Get("/khil", ws.ConnectKhil(ctx, lb))
}
