package ws

import (
	"context"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type Khil struct {
	User
}

func ConnectKhil(ctx context.Context, lb *Lobby) fiber.Handler {
	return websocket.New(func(conn *websocket.Conn) {
		uid := conn.Locals("uid")

		lb.Khil = &User{
			ID:       uid.(string),
			Conn:     conn,
			Receiver: make(chan *message, lb.PlayersAmount),
			Sender:   lb.KhilBC,
		}

		go func() {
			if err := lb.Khil.write(lb); err != nil {
				log.Info("Khil write error:", err)
				return
			}
		}()
		go lb.Khil.read(lb)

		log.Info("khil connected")

		select {
		case <-ctx.Done():
			return
		}
	})
}
