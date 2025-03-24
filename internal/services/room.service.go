package services

import (
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"sigo/internal/lib"
	"sigo/internal/models"
	"sort"
	"strconv"
	"time"
)

var (
	errTemp                = errors.New("error temp")
	errIdxOutOfRange       = errors.New("index out of range")
	errQuestionAlreadyUsed = errors.New("question already used")
)

const (
	typeConnected    = "connected"
	typeDisconnected = "disconnected"

	stageWaitForPlayers    = "wait_for_start"
	stageQuestionSelection = "question_selection"
	stageQuestion          = "question"
	stageButton            = "button"
	stageResponse          = "response"
	stageAnswer            = "answer"

	typeError          = "error"
	typeStart          = "start"
	typeTimeOut        = "send_time_out"
	typeNext           = "next"
	typeAcceptAnswer   = "accept_answer"
	typeDenyAnswer     = "deny_answer"
	typeSetStage       = "set_stage"
	typeModifyScore    = "modify_score"
	typePressButton    = "press_button"
	typeSetChooser     = "set_chooser"
	typeSlide          = "slide"
	typeQuestionSelect = "select_question"
	typeUsedQuestions  = "used_questions"
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
	log.Infof("Start listening room %v", room.Id())

	questionTimeOut := time.Second * time.Duration(room.Config().QuestionTime)

	room.Statement().SetStage(stageWaitForPlayers)
	room.Statement().SetUsedQuestions(make([][]bool, len(room.Pack().Rounds[0].Themes)))
	for i := range room.Statement().UsedQuestions() {
		room.Statement().UsedQuestions()[i] = make([]bool, len(room.Pack().Rounds[0].Themes[i].Questions))
	}

	timeToThinkAfterPressingTheButtonTimeOut := time.Second * time.Duration(room.Config().TimeToThinkAfterPressingTheButton)
	_ = timeToThinkAfterPressingTheButtonTimeOut
	log.Infof("%v", room.Statement())
	for {
		select {
		case request := <-*room.Receiver():
			switch request.Type {
			default:
				sendError(room, request.UID, errTemp)
			case typeConnected:
				response := lib.Response{
					UID:  request.UID,
					Type: request.Type,
				}
				sendForAllUsers(response, room)
			case typeDisconnected:
				response := lib.Response{
					UID:  request.UID,
					Type: request.Type,
				}
				sendForAllUsers(response, room)
			case typeStart:
				err := startGame(room, request)
				if err != nil {
					sendError(room, request.UID, err)
					continue
				}

				room.Statement().SetStage(stageQuestionSelection)
				response := lib.Response{
					Type: typeSetStage,
					Data: lib.Data{
						Content: stageQuestionSelection,
					},
				}
				sendForAllUsers(response, room)

				response = lib.Response{
					Type: typeSetChooser,
					Data: lib.Data{
						Content: strconv.FormatInt(room.Statement().ChooserId(), 10),
					},
				}
				sendForAllUsers(response, room)

				response = lib.Response{
					Type: typeUsedQuestions,
					Data: lib.Data{
						UsedQuestions: room.Statement().UsedQuestions(),
					},
				}
				sendForAllUsers(response, room)
			case typeQuestionSelect:
				err, question := selectQuestion(room, request)
				if err != nil {
					sendError(room, request.UID, err)
					continue
				}
				themeIdx := request.Data.ThemeIdx
				questionIdx := request.Data.QuestionIdx

				room.Statement().SetQuestion(question)
				room.Statement().SetStage(stageQuestion)

				response := lib.Response{
					Type: typeQuestionSelect,
					Data: lib.Data{
						ThemeIdx:    themeIdx,
						QuestionIdx: questionIdx,
					},
				}
				sendForAllUsers(response, room)

				slideIdx := room.Statement().SlideIdx()
				slide := question.QuestionSlides[slideIdx]

				response = lib.Response{
					Type: typeSlide,
					Data: lib.Data{
						Content:     *slide.Content,
						ContentType: *slide.ContentType,
					},
				}
				sendForAllUsers(response, room)

				room.Statement().SetUsedQuestion(themeIdx, questionIdx)
			case typeNext:
				err, slide := getNextSlide(room, request)
				if errors.Is(err, errIdxOutOfRange) {
					switch room.Statement().Stage() {
					case stageQuestion:
						room.Statement().SetSlideIdx(-1) // to use typeNext to show answer

						room.Statement().SetStage(stageButton)
						response := lib.Response{
							Type: typeSetStage,
							Data: lib.Data{
								Content: stageButton,
								Time:    time.Now().Add(questionTimeOut).Unix(),
							},
						}
						sendForAllUsers(response, room)

						go func() {
							time.Sleep(questionTimeOut)
							*room.TimeoutChannel() <- nil
						}()
						continue
					case stageAnswer:
						room.Statement().SetSlideIdx(0)

						room.Statement().SetStage(stageQuestionSelection)
						response := lib.Response{
							Type: typeSetStage,
							Data: lib.Data{
								Content:       stageQuestionSelection,
								UsedQuestions: room.Statement().UsedQuestions(),
							},
						}
						sendForAllUsers(response, room)
						continue
					}
				}
				if err != nil {
					sendError(room, request.UID, err)
					continue
				}

				response := lib.Response{
					Type: typeSlide,
					Data: lib.Data{
						Content:     *slide.Content,
						ContentType: *slide.ContentType,
					},
				}
				sendForAllUsers(response, room)
			}
		case <-*room.TimeoutChannel():
			response := lib.Response{
				Type: typeTimeOut,
			}
			sendForAllUsers(response, room)

			room.Statement().SetStage(stageAnswer)
			response = lib.Response{
				Type: typeSetStage,
				Data: lib.Data{
					Content: stageAnswer,
				},
			}
			sendForAllUsers(response, room)
		}
	}
}

func getNextSlide(room *models.Room, request lib.Request) (error, *lib.Slide) {
	if request.UID == room.Owner().Id() && room.Statement().Stage() == stageQuestion || room.Statement().Stage() == stageAnswer {
		slideIdx := room.Statement().SlideIdx()
		slides := make([]*lib.Slide, 0)

		switch room.Statement().Stage() {
		case stageQuestion:
			slides = room.Statement().Question().QuestionSlides
		case stageAnswer:
			slides = room.Statement().Question().AnswerSlides
		}

		if slideIdx+1 < len(slides) {
			room.Statement().SetSlideIdx(room.Statement().SlideIdx() + 1)
			return nil, slides[room.Statement().SlideIdx()]
		}

		return errIdxOutOfRange, nil
	}
	return errTemp, nil
}

func sendError(room *models.Room, requestUID int64, err error) {
	response := lib.Response{
		UID:  0,
		Type: typeError,
		Data: lib.Data{
			Content: err.Error(),
		},
	}
	*room.Players()[requestUID].Receiver() <- response
}

func selectQuestion(room *models.Room, request lib.Request) (error, *lib.Question) {
	if request.UID == room.Owner().Id() || request.UID == room.Statement().ChooserId() && room.Statement().Stage() == stageQuestionSelection {
		themeIdx := request.Data.ThemeIdx
		questionIdx := request.Data.QuestionIdx
		if questionIdx <= len(room.Statement().UsedQuestions()[themeIdx]) && themeIdx <= len(room.Statement().UsedQuestions()) {
			if room.Statement().UsedQuestions()[themeIdx][questionIdx] == false {
				return nil, room.Pack().Rounds[room.Statement().RoundIdx()].Themes[themeIdx].Questions[questionIdx]
			} else {
				return errQuestionAlreadyUsed, nil
			}
		}
	}
	return errTemp, nil
}

func sendForAllUsers(response lib.Response, room *models.Room) {
	for _, user := range room.Players() {
		*user.Receiver() <- response
	}
}

func startGame(room *models.Room, request lib.Request) error {
	if request.UID == room.Owner().Id() && room.Statement().Stage() == stageWaitForPlayers {
		if len(room.Players()) < 2 {
			return errTemp
		}
		var chooser *models.User
		for _, chooser = range room.Players() {
			if chooser.Id() != room.Owner().Id() {
				room.Statement().SetChooserId(chooser.Id())
				break
			}
		}
		return nil
	}

	return errTemp
}
