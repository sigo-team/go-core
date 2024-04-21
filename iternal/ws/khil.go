package ws

import (
	"context"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"math/rand"
	"sync"
)

type Khil struct {
	User
}

func ConnectKhil(ctx context.Context, lb *Lobby) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		lb.Khil = &User{
			ID:      rand.Int(),
			Conn:    c,
			Message: make(chan *message),
		}

		mu := &sync.Mutex{}
		go lb.Khil.write(lb, mu)
		go lb.Khil.read(lb, mu)

		log.Info("khil connected")

		select {
		case <-ctx.Done():
			return
		}
	})
}
