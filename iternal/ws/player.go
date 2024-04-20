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
	slg "sigo/iternal/.logger"
	"sigo/iternal/game"
)

type Player struct {
	User
	Name  string
	Score int
	// TODO: id
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

		msg := &message{
			User:    &player.User,
			Content: []byte(fmt.Sprintf("A new user has joined the lobby: %s", player.Name)),
		}
		lb.Broadcast <- msg

		go player.write(lb)
		go player.read(lb)
		/*
			go player.writeMessage(lb)
			go player.readMessage(lb)*/

		player.sendSiPackage(lb.SiPck)
		player.Conn.WriteMessage(1, []byte(fmt.Sprintf("%s", lb.Players)))

		select {
		case <-ctx.Done():
			lb.Unregister <- player
			return
		}

	})
}

func (player *Player) writeMessage(lb *Lobby) {
	defer func() {
		lb.Unregister <- player
		player.Conn.Close()
	}()

	for {
		select {
		case msg := <-player.Message:
			if err := player.Conn.WriteMessage(1, msg.Content); err != nil {
				log.Errorf("Error writing message: %v", err)
			}
		}
	}
}

func (player *Player) readMessage(lb *Lobby) {
	defer player.Conn.Close()

	for {
		_, data, err := player.Conn.ReadMessage()
		if err != nil {
			return
		}

		msg := &message{
			User:    &player.User,
			Content: data,
		}

		lb.Broadcast <- msg
	}
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

	if err := player.Conn.WriteJSON(siPckHeaders); err != nil {
		slg.Err(err)
	}
}
