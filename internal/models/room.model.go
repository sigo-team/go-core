package models

import (
	"encoding/json"
	"errors"
	"os"
	"sigo/internal/lib"
)

type RoomConfig struct {
	Public                                bool `form:"public"`
	RoundTime                             uint `form:"round_time"`
	TimeToThinkAfterPressingTheButton     uint `form:"time_to_think_after_pressing_the_button"`
	QuestionTime                          int  `form:"question_time"`
	TimeToThinkOnASpecialQuestion         uint `form:"time_to_think_on_a_special_question"`
	TimeToThinkAtTheFinal                 uint `form:"time_to_think_at_the_final"`
	EndTheQuestionIfTheAnswerIsCorrect    bool `form:"end_the_question_if_the_answer_is_correct"`
	GameWithFalseStarts                   bool `form:"game_with_false_starts"`
	MultimediaWithFalseStarts             bool `form:"multimedia_with_false_starts"`
	PlaySpecialQuestions                  bool `form:"play_special_questions"`
	RollBackStatisticsWhenTakingAStepBack bool `form:"roll_back_statistics_when_taking_a_step_back"`
	AutomaticGamePlay                     bool `form:"automatic_game_play"`
	DeductPointsForAnIncorrectAnswer      bool `form:"deduct_points_for_an_incorrect_answer"`
	ShowPlayersWhoHaveLostTheButton       bool `form:"show_players_who_have_lost_the_button"`
}

type RoomOptions struct {
	Owner       *User
	PackageName string
	Config      RoomConfig
}

type Room struct {
	owner      *User
	players    map[int64]*User
	spectators map[int64]*User

	id       int64
	packName string
	pack     lib.Pack

	scoreTab       map[int64]int
	statement      *Statement
	receiver       chan lib.Request
	sender         chan lib.Response
	timeoutChannel chan interface{}

	config RoomConfig
}

func (r *Room) MarshalJSON() ([]byte, error) {
	players := make(map[int64]string, len(r.players))
	for _, user := range r.players {
		players[user.id] = user.Name()
	}
	owner := struct {
		Uid  int64  `json:"uid"`
		Name string `json:"name"`
	}{
		Uid:  r.owner.id,
		Name: r.owner.name,
	}

	return json.Marshal(struct {
		Owner struct {
			Uid  int64  `json:"uid"`
			Name string `json:"name"`
		} `json:"owner"`
		Players       map[int64]string `json:"players"`
		PlayersAmount int              `json:"players_amount"`
		Id            int64            `json:"id"`
		PackageName   string           `json:"package_name"`
		IsPublic      bool             `json:"is_public"`
	}{
		Owner:         owner,
		Players:       players,
		PlayersAmount: len(r.players),
		Id:            r.Id(),
		PackageName:   r.packName,
		IsPublic:      r.config.Public,
	})
}

func (r *Room) Id() int64 {
	if r.id == 0 {
		panic("using room id before mounting it to the service")
	}
	return r.id
}

func (r *Room) Mount(id int64) {
	if r.id != 0 {
		panic("room is already mounted")
	}
	r.id = id
}

func validateRoomOptions(options RoomOptions) error {
	if options.Owner == nil {
		return errors.New("")
	}
	if len(options.PackageName) == 0 {
		return errors.New("")
	}
	return nil
}

func NewRoom(options RoomOptions) (*Room, error) {
	err := validateRoomOptions(options)
	if err != nil {
		return nil, err
	}
	content, err := os.ReadFile("./" + options.PackageName + "/content.json")
	if err != nil {
		return nil, err
	}
	pck := new(lib.Pack)
	err = json.Unmarshal(content, pck)
	if err != nil {
		return nil, err
	}

	return &Room{
		// FIXME: magic number
		owner:          options.Owner,
		players:        make(map[int64]*User, 0),
		spectators:     make(map[int64]*User, 0),
		packName:       options.PackageName,
		pack:           *pck,
		scoreTab:       make(map[int64]int, 0),
		receiver:       make(chan lib.Request),
		sender:         make(chan lib.Response),
		timeoutChannel: make(chan interface{}),
		config:         options.Config,
		statement:      new(Statement),
	}, nil
}

func (r *Room) JoinPlayer(user *User) {
	r.players[user.Id()] = user
}

func (r *Room) DisjoinPlayer(user *User) {
	delete(r.players, user.Id())
}

func (r *Room) ModifyScore(uid int64, delta int) {
	r.scoreTab[uid] += delta
}

func (r *Room) Owner() *User {
	return r.owner
}

func (r Room) Players() map[int64]*User {
	return r.players
}

func (r *Room) Receiver() *chan lib.Request {
	return &r.receiver
}

func (r *Room) Statement() *Statement {
	return r.statement
}

func (r *Room) PackageName() string {
	return r.packName
}

func (r *Room) Pack() lib.Pack {
	return r.pack
}

func (r *Room) Config() RoomConfig {
	return r.config
}

func (r *Room) TimeoutChannel() *chan interface{} {
	return &r.timeoutChannel
}
