package ws

import (
	"context"
	"encoding/json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/loremipsum.v1"
	"sigo/internal/lib"
)

type Player struct {
	User
	Name  string
	Score int
}

// todo: write and read to User

func ConnectPlayerHandler(ctx context.Context, lb *Lobby) fiber.Handler {
	return websocket.New(func(conn *websocket.Conn) {
		uid := conn.Locals("uid")

		player := &Player{
			User: User{
				ID:       uid.(string),
				Conn:     conn,
				Receiver: make(chan *message, lb.PlayersAmount),
				Sender:   make(chan *message, lb.PlayersAmount),
			},
			Name:  loremipsum.New().Word(),
			Score: 0,
		}

		lb.Register <- player

		go func() {
			if err := player.write(lb); err != nil {
				lb.Unregister <- player
				return
			}
		}()
		go player.read(lb)

		select {
		case <-ctx.Done():
			lb.Unregister <- player
			return
		}

	})
}

func (player *Player) sendSiPackage(p *lib.Package) {
	marshal, err := json.Marshal(p)
	if err != nil {
		return
	}

	siPckHeaders := new(lib.Package)
	err = json.Unmarshal(marshal, &siPckHeaders)
	if err != nil {
		return
	}

	for _, round := range siPckHeaders.Rounds {
		for _, theme := range round.Themes {
			for _, question := range theme.Questions {
				question.Type = nil
				question.AnswerSlides = nil
				question.QuestionSlides = nil
				question.PriceMax = nil
				question.PriceStep = nil
			}
		}
	}

	bytes, err := json.Marshal(siPckHeaders)
	if err != nil {
		return
	}
	player.Receiver <- &message{
		UserID:  "",
		Content: bytes,
	}
}
