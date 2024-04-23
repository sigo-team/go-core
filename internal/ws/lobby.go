package ws

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"sigo/internal/api"
	"sigo/internal/lib"
	"time"
)

type Lobby struct {
	Statement     *Statement
	SiPck         *lib.Package
	Khil          *User
	Chooser       *Player
	Respondent    *Player
	PlayersAmount uint
	Players       map[string]*Player
	ButtonBC      chan *message
	ChooserBC     chan *message
	KhilBC        chan *message

	Register   chan *Player
	Unregister chan *Player
}

type Statement struct {
	Status        string
	RoundIndex    uint
	Question      *lib.Question
	SlideIndex    int
	UsedQuestions [][]bool
}

func NewLobby(siPck *lib.Package, playersAmount uint) *Lobby {
	// TODO: config
	lb := &Lobby{
		Statement: &Statement{
			Status:     "waitForStart",
			RoundIndex: 0,
			SlideIndex: 0,
		},
		SiPck:         siPck,
		Khil:          new(User),
		Chooser:       new(Player),
		Respondent:    new(Player),
		PlayersAmount: playersAmount,
		Players:       make(map[string]*Player),
		ButtonBC:      make(chan *message),
		ChooserBC:     make(chan *message),
		KhilBC:        make(chan *message),
		Register:      make(chan *Player, playersAmount),
		Unregister:    make(chan *Player, playersAmount),
	}

	lb.Statement.UsedQuestions = make([][]bool, len(lb.SiPck.Rounds[0].Themes))
	for i := range lb.Statement.UsedQuestions {
		lb.Statement.UsedQuestions[i] = make([]bool, len(lb.SiPck.Rounds[0].Themes[i].Questions))
	}

	return lb
}

func (lb *Lobby) RunLobby(ctx context.Context) {
	buttoningTO := 5 * time.Second
	answeringTO := 5 * time.Second

	for {
		select {
		// Broadcast
		case msg := <-lb.KhilBC:
			{
				log.Debugf("%s", lb.Statement.Status)
				response := api.ReadResponse(msg.Content)

				switch lb.Statement.Status {
				case "waitForStart":
					{
						switch response.Type {
						case "start":
							{
								if len(lb.Players) == 0 {
									log.Error("Can't start lobby, no players left")
									request := api.Request{
										Type: "error",
										Data: api.Data{
											Content: "Can't start lobby, no players left",
										},
									}.Marshall()

									lb.Khil.Receiver <- &message{
										UserID:  "",
										Content: request,
									}

									break
								}

								for _, player := range lb.Players {
									lb.SetChooser(player)
									break
								}
								lb.Statement.Status = "choosing"
								log.Infof("Game started, chooser: %s", lb.Chooser.ID)
							}
						}
					}
				case "question":
					{
						switch response.Type {
						case "nextSlide":
							lb.Statement.SlideIndex++
							if lb.Statement.SlideIndex >= len(lb.Statement.Question.QuestionSlides) {
								lb.Statement.SlideIndex = 0
								lb.Statement.Status = "buttoning"
								request := api.Request{
									Type: "setButtoning",
								}.Marshall()

								lb.SendAll(&message{
									UserID:  "",
									Content: request,
								})

								log.Info("Start buttoning")
								break
							}
							lb.SendSlide()
						}
					}

					// TODO: answering
				case "answering":
					switch response.Type {
					case "accepted":
						{
							scoreChanges := +*lb.Statement.Question.PriceMin
							lb.Respondent.Score += scoreChanges
							lb.SetChooser(lb.Respondent)
							lb.Respondent = nil
							request := api.Request{
								Type: "changeScore",
								Data: api.Data{
									PlayerId:     lb.Respondent.ID,
									ScoreChanges: scoreChanges,
								},
							}.Marshall()

							lb.SendAll(&message{
								UserID:  "",
								Content: request,
							})

						}
					case "denied":
						{
							scoreChanges := -*lb.Statement.Question.PriceMin
							lb.Respondent.Score += scoreChanges
							lb.Respondent = nil
							request := api.Request{
								Type: "changeScore",
								Data: api.Data{
									PlayerId:     lb.Respondent.ID,
									ScoreChanges: scoreChanges,
								},
							}.Marshall()

							lb.SendAll(&message{
								UserID:  "",
								Content: request,
							})

						}

					}
				}
			}
		case msg := <-lb.ChooserBC:
			{
				response := api.ReadResponse(msg.Content)

				switch lb.Statement.Status {
				case "choosing":
					if response.Type == "chooseQuestion" {
						roundIdx := lb.Statement.RoundIndex

						themeIdx := response.Data.ThemeIndex
						questionIdx := response.Data.QuestionIndex

						lb.Statement.Question = lb.SiPck.Rounds[roundIdx].Themes[themeIdx].Questions[questionIdx]

						request := api.Request{
							Type: "chooseQuestion",
							Data: api.Data{
								ThemeIndex:    themeIdx,
								QuestionIndex: questionIdx,
							},
						}.Marshall()

						lb.SendAll(&message{
							UserID:  lb.Chooser.ID,
							Content: request,
						})

						lb.Statement.Status = "question"
						lb.SendSlide()
						log.Info("Started question")
					}
				}
			}

		case msg := <-lb.ButtonBC:
			{
				response := api.ReadResponse(msg.Content)
				_ = response

				switch lb.Statement.Status {
				case "buttoning":
					lb.SendAll(msg)
					lb.Respondent = lb.Players[msg.UserID]
					lb.Statement.Status = "answering"
				}
			}

			//TODO: check timeouts work
		case <-time.After(buttoningTO):
			{
				if lb.Statement.Status == "buttoning" {
					request := api.Request{
						Type: "timeout",
						Data: api.Data{
							Content: "buttoning",
						},
					}.Marshall()
					lb.SendAll(&message{
						UserID:  "",
						Content: request,
					})
					lb.Statement.Status = "choosing"
					log.Info("Button timeout. Started choosing")
				}
			}

		case <-time.After(answeringTO):
			{
				if lb.Statement.Status == "answering" {
					request := api.Request{
						Type: "timeout",
						Data: api.Data{
							Content: "answering",
						},
					}.Marshall()
					lb.SendAll(&message{
						UserID:  "",
						Content: request,
					})
				}
			}

		// Other
		case player := <-lb.Register:
			lb.AddPlayer(player)
		case player := <-lb.Unregister:
			lb.RemovePlayer(player)
		case <-ctx.Done():
			return
		}
	}
}

func (lb *Lobby) AddPlayer(player *Player) {
	lb.Players[player.ID] = player
	player.sendSiPackage(lb.SiPck)
	content := api.Request{
		Type: "playerConnect",
		Data: api.Data{PlayerId: player.ID},
	}.Marshall()

	lb.SendAll(&message{
		UserID:  player.ID,
		Content: content,
	})

	log.Infof("Player %s joined the game", player.Name)
}

func (lb *Lobby) RemovePlayer(player *Player) {
	delete(lb.Players, player.ID)
	content := api.Request{
		Type: "playerDisconnect",
		Data: api.Data{PlayerId: player.ID},
	}.Marshall()

	lb.SendAll(&message{
		UserID:  player.ID,
		Content: content,
	})
	log.Infof("Player %s left the game", player.Name)
}

func (lb Lobby) SendAll(msg *message) {
	for _, player := range lb.Players {
		player.Receiver <- msg
	}
	lb.Khil.Receiver <- msg
}

func (lb Lobby) SendSlide() {
	slideIdx := lb.Statement.SlideIndex

	request := api.Request{
		Type: "sendSlide",
		Data: api.Data{
			Type:    *lb.Statement.Question.QuestionSlides[slideIdx].ContentType,
			Content: *lb.Statement.Question.QuestionSlides[slideIdx].Content,
		},
	}.Marshall()

	lb.SendAll(&message{
		UserID:  "",
		Content: request,
	})
}
