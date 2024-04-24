package transport

import (
	"context"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"sigo/internal/services"
)

func ConnectPlayer(ctx context.Context, us *services.UserService, room *services.Room) fiber.Handler {
	return websocket.New(func(conn *websocket.Conn) {
		defer conn.Close()

		player := services.NewPlayer(us)

		go read(conn, &player.User)
		go write(conn, &player.User)

		//todo: register

		select {
		case <-ctx.Done():
		}
	})
}
