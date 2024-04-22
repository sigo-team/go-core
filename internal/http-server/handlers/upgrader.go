package handlers

import (
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"math/rand"
	"strconv"
)

func UpgradeHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(ctx) {
			uid := GetUserId()
			log.Info(fmt.Sprintf("Successfully upgrade required %s", uid))
			ctx.Locals("uid", uid)
			return ctx.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}

func GetUserId() string {
	return strconv.Itoa(rand.Int())
}
