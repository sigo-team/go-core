package transport

import (
	"context"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"sigo/internal/services"
)

func ConnectPlayer(us *services.UserService, ctx context.Context, room *services.Room) fiber.Handler {
	return websocket.New(func(conn *websocket.Conn) {
		defer conn.Close()

		player := services.NewPlayer(us)

		go transitIn(conn, &player.User)
		go transitOut(conn, &player.User)

		//todo: register, closing

		select {
		case <-ctx.Done():
		}
	})
}
