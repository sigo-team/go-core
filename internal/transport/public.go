package transport

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"sigo/internal/controllers"
)

func PublicRoutes(closingCtx context.Context, app *fiber.App, roomController *controllers.RoomController) {
	route := app.Group("/api/v1")

	route.Get("/room", roomController.GetRooms)
	route.Post("/room", roomController.CreateRoom)

	route.Get("/media/:packageName/:fileName", sendMedia)

	route.Use("/ws", UpgradeMiddleware)
	route.Get("/ws", controllers.Handler(closingCtx, roomController))
}

func sendMedia(ctx *fiber.Ctx) error {
	packageName := ctx.Params("packageName")
	fileName := ctx.Params("fileName")
	return ctx.SendFile("./" + packageName + "/media/" + fileName)
}
