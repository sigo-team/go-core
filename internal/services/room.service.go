package services

import (
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"sigo/internal/lib"
	"sigo/internal/models"
	"sort"
	"time"
)

const (
	answerSlide   = 0
	questionSlide = 1

	waitForStartStage       = "waitForStart"
	questionSelectionStage  = "questionSelection"
	showingQuestionStage    = "showingQuestion"
	waitForPressButtonStage = "waitForPressButton"
	givingAnswerStage       = "givingAnswer"
	showingAnswerStage      = "showingAnswer"

	infoType           = "info"
	errorType          = "error"
	startType          = "start"
	timeOutType        = "sendTimeOut"
	nextType           = "next"
	pastType           = "past"
	acceptAnswerType   = "acceptAnswer"
	denyAnswerType     = "cancelAnswer"
	setStageType       = "setStage"
	modifyScoreType    = "modifyScore"
	pressButtonType    = "pressButton"
	setChooserType     = "setChooser"
	slideType          = "slide"
	questionSelectType = "questionSelect"
)

type RoomService struct {
	identifierManager *lib.IdentifierManager
	rooms             map[int64]*models.Room
}

type RoomServiceOptions struct {
	IdentifierManager *lib.IdentifierManager
}

func validateRoomServiceOptions(options RoomServiceOptions) error {
	return nil
}

func NewRoomService(options RoomServiceOptions) *RoomService {
	err := validateRoomServiceOptions(options)
	if err != nil {
		panic(err)
	}
	return &RoomService{identifierManager: options.IdentifierManager, rooms: make(map[int64]*models.Room)}
}

func (r *RoomService) CreteRoom(options models.RoomOptions) (*models.Room, error) {
	room, err := models.NewRoom(options)
	if err != nil {
		return nil, err
	}
	room.Mount(r.identifierManager.NewID())
	r.rooms[room.Id()] = room
	return room, nil
}

func (r *RoomService) ReadRoom(roomId int64) (*models.Room, error) {
	room, ok := r.rooms[roomId]
	if !ok {
		return nil, errors.New("")
	}
	return room, nil
}

func partitioned(p int, n int, slice []int64) []int64 {
	if n < 0 {
		return []int64{}
	}
	start := (n - 1) * p
	end := n * p
	return slice[start:min(end, len(slice))]
}

func (r *RoomService) GetRoomsAmount() int {
	return len(r.rooms)
}

func (r *RoomService) ReadRooms(page int) ([]*models.Room, error) {
	ids := make([]int64, 0)
	for _, room := range r.rooms {
		ids = append(ids, room.Id())
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
	part := partitioned(8, page, ids)
	rooms := make([]*models.Room, 0)
	for _, id := range part {
		rooms = append(rooms, r.rooms[id])
	}
	return rooms, nil
}

func Listening(room *models.Room) {
	log.Info("Start listening room")

	*room.Statement() = models.Statement{
		Stage:        waitForStartStage,
		RoundIdx:     0,
		Question:     nil,
		SlideIdx:     0,
		AnswerableID: 0,
	}

	questionTimeOut := time.Second * time.Duration(room.Config().QuestionTime)
	timeToThinkAfterPressingTheButtonTimeOut := time.Second * time.Duration(room.Config().TimeToThinkAfterPressingTheButton)

	for {
		switch room.Statement().Stage {
		case waitForStartStage:
			select {
			case request := <-*room.Owner().Sender():
				waitForStartOwnerChecker(room, &request)
			}
		case questionSelectionStage:
			select {
			case request := <-*room.ChooserBC():
				questionSelectionChooserChecker(room, &request)
			}
		case showingQuestionStage:
			select {
			case request := <-*room.Owner().Sender():
				showingQuestionOwnerChecker(room, &request)
			}
		case waitForPressButtonStage:
			select {
			case request := <-*room.ButtonBC():
				buttonChecker(room, &request)
			case <-time.After(questionTimeOut):
				sendTimeOut(room, "questionTimeOut")
			}
		case showingAnswerStage:
			select {
			case request := <-*room.Owner().Sender():
				showingAnswerOwnerChecker(room, &request)
			}
		case givingAnswerStage:
			select {
			case request := <-*room.Owner().Sender():
				givingAnswerOwnerChecker(room, &request)
			case <-time.After(timeToThinkAfterPressingTheButtonTimeOut):
				sendTimeOut(room, "timeToThinkAfterPressingTheButtonTimeOut")
			}
		}
	}
}

func waitForStartOwnerChecker(room *models.Room, request *lib.Request) {
	switch request.Type {
	case startType:
		if len(room.Players()) == 0 {
			log.Errorf("Room %v has no players, cannot start", room.Id())
			*room.Owner().Receiver() <- lib.Response{
				UID:  0,
				Type: errorType,
				Data: lib.Data{
					Content: "Room has no players",
				},
			}
			return
		}
		for _, user := range room.Players() {
			user.SetSender(room.ChooserBC())
			setChooser(room, user)
			break
		}
		room.Statement().Stage = questionSelectionStage
	}
}

func questionSelectionChooserChecker(room *models.Room, request *lib.Request) {
	switch request.Type {
	case questionSelectionStage:
		roundIdx := room.Statement().RoundIdx
		themeIdx := request.Data.ThemeIndex
		questionIdx := request.Data.QuestionIndex

		response := lib.Response{
			UID:  0,
			Type: questionSelectType,
			Data: lib.Data{
				ThemeIndex:    themeIdx,
				QuestionIndex: questionIdx,
			},
		}

		sendForAll(room, response)

		log.Infof("Room %v has choose: %d %d", room.Id(), themeIdx, questionIdx)
		for _, user := range room.Players() {
			user.SetSender(room.ButtonBC())
		}
		room.Statement().Stage = showingAnswerStage
		room.Statement().Question = room.Pack().Rounds[roundIdx].Themes[themeIdx].Questions[questionIdx]

		*room.Owner().Receiver() <- lib.Response{
			Type: showingAnswerStage,
			Data: lib.Data{
				Question: *room.Statement().Question,
			},
		}

		slide := *room.Statement().Question.QuestionSlides[room.Statement().SlideIdx]
		sendSlide(room, slide)
	}
}

func showingQuestionOwnerChecker(room *models.Room, request *lib.Request) {
	switch request.Type {
	case nextType:
		room.Statement().SlideIdx++
		if len(room.Statement().Question.QuestionSlides) <= room.Statement().SlideIdx {
			room.Statement().Stage = waitForPressButtonStage
			sendStage(room, waitForPressButtonStage)
			room.Statement().SlideIdx = 0
			return
		}

		slide := *room.Statement().Question.QuestionSlides[room.Statement().SlideIdx]
		sendSlide(room, slide)
	}
}

func buttonChecker(room *models.Room, request *lib.Request) {
	switch request.Type {
	case pressButtonType:
		log.Infof("In room %d user %d pressed the button", room.Id(), request.UID)
		room.Statement().Stage = givingAnswerStage
		room.Statement().AnswerableID = request.UID

		sendStage(room, givingAnswerStage)
	}
}

func showingAnswerOwnerChecker(room *models.Room, request *lib.Request) {
	switch request.Type {
	case nextType:
		room.Statement().SlideIdx++
		if room.Statement().SlideIdx >= len(room.Statement().Question.AnswerSlides) {
			room.Statement().SlideIdx = 0

			room.Statement().Stage = questionSelectionStage
			sendStage(room, questionSelectionStage)
			return
		}

		slide := *room.Statement().Question.AnswerSlides[room.Statement().SlideIdx]
		sendSlide(room, slide)
	}
}

func givingAnswerOwnerChecker(room *models.Room, request *lib.Request) {
	var delta int
	switch request.Type {
	case acceptAnswerType:
		delta = *room.Statement().Question.PriceMin
	case denyAnswerType:
		delta = -*room.Statement().Question.PriceMin
	}
	if request.Type == acceptAnswerType || request.Type == denyAnswerType {
		room.ModifyScore(request.UID, delta)
		response := lib.Response{
			UID:  0,
			Type: modifyScoreType,
			Data: lib.Data{
				ScoreChanges: delta,
				UID:          request.UID,
			},
		}
		sendForAll(room, response)

		room.Statement().Stage = showingAnswerStage
	}
}

func setChooser(room *models.Room, user *models.User) {
	response := lib.Response{
		UID:  0,
		Type: setChooserType,
		Data: lib.Data{
			UID: user.Id(),
		},
	}

	for _, u := range room.Players() {
		*u.Receiver() <- response
	}
	*room.Owner().Receiver() <- response
	log.Infof("Set Chooser %v %v", room.Id(), user.Id())
}

func sendSlide(room *models.Room, slide lib.Slide) {
	content := *slide.Content
	contentType := *slide.ContentType

	response := lib.Response{
		UID:  0,
		Type: slideType,
		Data: lib.Data{
			Content:     content,
			ContentType: contentType,
		},
	}
	sendForAll(room, response)
}

func sendTimeOut(room *models.Room, content string) {
	log.Infof("Room %d waitButton timeout", room.Id())
	response := lib.Response{
		Type: timeOutType,
		Data: lib.Data{
			Content: content,
		},
	}
	sendForAll(room, response)

	room.Statement().Stage = showingAnswerStage
}

func sendStage(room *models.Room, stage string) {
	response := lib.Response{
		UID:  0,
		Type: setStageType,
		Data: lib.Data{
			ContentType: stage,
		},
	}

	sendForAll(room, response)
	log.Infof("Room %d set stage %s", room.Id(), stage)
}

func sendForAll(room *models.Room, response lib.Response) {
	for _, user := range room.Players() {
		*user.Receiver() <- response
	}
	*room.Owner().Receiver() <- response
}
