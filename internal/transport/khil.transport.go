package transport

import (
	"context"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"sigo/internal/services"
)

func ConnectKhil(ctx context.Context, us *services.UserService, room *services.Room) fiber.Handler {
	return websocket.New(func(conn *websocket.Conn) {
		defer conn.Close()

		khil := services.NewKhil(us)

		go read(conn, &khil.User)
		go write(conn, &khil.User)

		//todo: register

		select {
		case <-ctx.Done():
		}
	})
}
