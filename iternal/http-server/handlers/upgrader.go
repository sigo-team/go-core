package handlers

import (
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func UpgradeHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(ctx) {
			log.Info(fmt.Sprintf("Successfully upgrade required %s", ctx.IP()))
			return ctx.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}
