package transport

import (
	"context"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"sigo/internal/controllers"
)

func ConnectKhil(ctx context.Context, r *controllers.RoomHandlers) fiber.Handler {
	return websocket.New(func(conn *websocket.Conn) {
		defer conn.Close()

		userID := conn.Locals(UserIDKey).(int64)

		go transitIn(conn, r.UserService.GetUser(userID))
		go transitOut(conn, r.UserService.GetUser(userID))

		select {
		case <-ctx.Done():
		}
	})
}
