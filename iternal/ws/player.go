package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gopkg.in/loremipsum.v1"
	"math/rand"
	"sigo/iternal/api"
	"sigo/iternal/game"
	"sync"
)

type Player struct {
	User
	Name  string
	Score int
}

// todo: write and read to User

func ConnectPlayerHandler(ctx context.Context, lb *Lobby) fiber.Handler {
	return websocket.New(func(conn *websocket.Conn) {

		player := &Player{
			User: User{
				ID:      rand.Int(),
				Conn:    conn,
				Message: make(chan *message, lb.PlayersAmount),
			},
			Name:  loremipsum.New().Word(),
			Score: 0,
		}

		lb.Register <- player

		request := api.Request{
			Type: "greeting",
			Data: api.Data{
				ThemeIndex:    0,
				QuestionIndex: 0,
				PlayerId:      player.ID,
				ScoreChanges:  0,
				ChooserID:     0,
			},
		}
		marshalled, err := json.Marshal(request)
		if err != nil {
			log.Errorf("Error marshalling request: %v", err)
		}

		msg := &message{
			User:    &player.User,
			Content: []byte(marshalled),
		}
		lb.Broadcast <- msg

		mu := &sync.Mutex{}

		go player.write(lb, mu)
		go player.read(lb, mu)

		player.sendSiPackage(lb.SiPck)
		player.Message <- &message{
			User:    nil,
			Content: []byte(fmt.Sprintf("%s", lb.Players)),
		}

		select {
		case <-ctx.Done():
			lb.Unregister <- player
			return
		}

	})
}

func (player *Player) sendSiPackage(p *game.Package) {
	marshal, err := json.Marshal(p)
	if err != nil {
		return
	}

	siPckHeaders := new(game.Package)
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
	player.Message <- &message{
		User:    nil,
		Content: bytes,
	}
}
