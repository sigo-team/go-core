package ws

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2/log"
	"sigo/iternal/api"
	"sigo/iternal/game"
)

type Lobby struct {
	Statement     Statement
	SiPck         *game.Package
	Khil          *User
	PlayersAmount uint
	Players       map[int]*Player
	Chooser       Player
	Broadcast     chan *message
	Register      chan *Player
	Unregister    chan *Player
}

type Statement struct {
	Status     string
	RoundIndex uint
	SlideIndex int
	Question   *game.Question
}

func NewLobby(siPck *game.Package, playersAmount uint) *Lobby {
	return &Lobby{
		Statement: Statement{
			Status:     "waitForStart",
			RoundIndex: 0,
			// TODO: починить этот кастыль
			SlideIndex: -1,
		},
		SiPck:         siPck,
		Khil:          new(User),
		PlayersAmount: playersAmount,
		Players:       make(map[int]*Player),
		Chooser:       Player{},
		Broadcast:     make(chan *message, playersAmount),
		Register:      make(chan *Player, playersAmount),
		Unregister:    make(chan *Player),
	}
}

func (lb *Lobby) RunLobby(ctx context.Context) {
	for {
		select {
		// Broadcast
		case msg := <-lb.Broadcast:
			response := new(api.Response)
			if err := json.Unmarshal(msg.Content, response); err != nil {
				log.Errorf("Error unmarshalling msg: %v", err)
			}

			if response.Type == "greeting" {
				for _, player := range lb.Players {
					player.Message <- msg
				}
				lb.Khil.Message <- msg
			}

			switch lb.Statement.Status {
			case "waitForStart":
				{
					if msg.User.ID == lb.Khil.ID {
						switch response.Type {
						case "start":
							{
								if len(lb.Players) != 0 {
									for _, player := range lb.Players {
										lb.Chooser = *player
										break
									}
									request := api.Request{
										Type: "choose",
										Data: api.Data{
											ChooserID: lb.Chooser.ID,
										},
									}

									marshal, err := json.Marshal(request)
									if err != nil {
										log.Errorf("Error marshalling req: %v", err)
									}

									for _, player := range lb.Players {
										player.Message <- &message{
											User:    lb.Khil,
											Content: marshal,
										}
									}
								} else {
									log.Errorf("Lobby has no players, cannot start")
								}
								log.Infof("chooserId: %v", lb.Chooser.ID)
								lb.Statement.Status = "choosing"
							}
						}
					}
				}
			case "choosing":
				{
					if msg.User.ID == lb.Chooser.ID {
						for _, player := range lb.Players {
							player.Message <- msg
						}
						lb.Khil.Message <- msg
						lb.Broadcast <- msg
						lb.Statement = Statement{
							Status:   "question",
							Question: lb.SiPck.Rounds[lb.Statement.RoundIndex].Themes[response.Data.ThemeIndex].Questions[response.Data.QuestionIndex],
						}
					}
				}
			case "question":
				{
					// TODO: закончить эту ебалу
					if msg.User.ID == lb.Khil.ID {
						if response.Type == "next" {
							log.Debugf("%v", lb.Statement.SlideIndex)
							if lb.Statement.SlideIndex < len(lb.Statement.Question.QuestionSlides) {
								for _, player := range lb.Players {
									player.Message <- &message{
										User:    nil,
										Content: []byte(*lb.Statement.Question.QuestionSlides[lb.Statement.SlideIndex].Content),
									}
								}
								lb.Statement.SlideIndex++
							} else {
								lb.Statement = Statement{
									Status:     "buttoning",
									Question:   nil,
									SlideIndex: 0,
								}
							}
						}
					}
				}
			case "buttoning":
				{
					// TODO: buttoning
					if response.Type == "pressButton" {
						for _, player := range lb.Players {
							player.Message <- msg
						}
					}
				}
			}

		// Other
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
