package ws

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"sigo/iternal/game"
)

type Lobby struct {
	SiPck         *game.Package
	Players       map[int]*Player
	PlayersAmount uint
	Khil          *Khil
	Broadcast     chan *message
	Register      chan *Player
	Unregister    chan *Player
}

func NewLobby(siPck *game.Package, playersAmount uint) *Lobby {
	return &Lobby{
		SiPck:         siPck,
		Players:       make(map[int]*Player),
		PlayersAmount: playersAmount,
		Khil:          new(Khil),
		Broadcast:     make(chan *message, playersAmount),
		Register:      make(chan *Player, playersAmount),
		Unregister:    make(chan *Player),
	}
}

func (lb *Lobby) WaitStart(ctx context.Context) {

}

func (lb *Lobby) RunLobby(ctx context.Context) {
	for {
		select {
		case msg := <-lb.Broadcast:
			switch msg.User.ID {
			case lb.Khil.ID:
				log.Info("khil send message")
				if string(msg.Content) == "next" {
					for _, player := range lb.Players {
						log.Info("player:", player.ID)
						player.Message <- &message{
							User:    &lb.Khil.User,
							Content: []byte("next"),
						}
					}
				}
			default:
				for _, player := range lb.Players {
					player.Message <- msg
				}
				lb.Khil.Message <- msg
			}
		case player := <-lb.Register:
			lb.AddPlayer(player)
		case player := <-lb.Unregister:
			lb.RemovePlayer(player)
			log.Info("got a unreg")
		case <-ctx.Done():
			return
		}
	}
}

func (lb *Lobby) AddPlayer(p *Player) {
	log.Infof("Added ID: %v", p.ID)
	lb.Players[p.ID] = p
}

func (lb *Lobby) RemovePlayer(player *Player) {
	delete(lb.Players, player.ID)
}
