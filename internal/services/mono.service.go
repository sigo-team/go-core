package services

import (
	"errors"
	"sigo/internal/lib"
	"strconv"
	"strings"
)

type MonoService struct {
	DB struct {
		idManager *lib.IdentifierManager
		Rooms     []Room
	}
}

type Room struct {
	ID                                    int64
	Public                                bool
	Players                               int
	MaxPlayers                            int
	PackageName                           string
	RoundTime                             int  `form:"round_time"`
	TimeToThinkAfterPressingTheButton     int  `form:"time_to_think_after_pressing_the_button"`
	QuestionTime                          int  `form:"question_time"`
	TimeToThinkOnASpecialQuestion         int  `form:"time_to_think_on_a_special_question"`
	TimeToThinkAtTheFinal                 int  `form:"time_to_think_at_the_final"`
	EndTheQuestionIfTheAnswerIsCorrect    bool `form:"end_the_question_if_the_answer_is_correct"`
	GameWithFalseStarts                   bool `form:"game_with_false_starts"`
	MultimediaWithFalseStarts             bool `form:"multimedia_with_false_starts"`
	PlaySpecialQuestions                  bool `form:"play_special_questions"`
	RollBackStatisticsWhenTakingAStepBack bool `form:"roll_back_statistics_when_taking_a_step_back"`
	AutomaticGamePlay                     bool `form:"automatic_game_play"`
	DeductPointsForAnIncorrectAnswer      bool `form:"deduct_points_for_an_incorrect_answer"`
	ShowPlayersWhoHaveLostTheButton       bool `form:"show_players_who_have_lost_the_button"`
}

func ValidateRoom(r Room) bool {
	if r.PackageName == "" {
		return false
	}
	if r.RoundTime <= 0 {
		return false
	}
	if r.TimeToThinkAfterPressingTheButton <= 0 {
		return false
	}
	if r.QuestionTime <= 0 {
		return false
	}
	if r.TimeToThinkOnASpecialQuestion <= 0 {
		return false
	}
	if r.TimeToThinkAtTheFinal <= 0 {
		return false
	}
	return true
}

func NewMonoService() *MonoService {
	return &MonoService{
		DB: struct {
			idManager *lib.IdentifierManager
			Rooms     []Room
		}{idManager: lib.NewIdentifierManager(), Rooms: make([]Room, 0)},
	}
}

var (
	OutOfRangeErr  = errors.New("out of range")
	TooShortKeyErr = errors.New("too short key")
	ValidationErr  = errors.New("room fields are not valid")
)

func (s *MonoService) CreateRoom(r Room) (Room, error) {
	r.ID = s.DB.idManager.NewID()
	r.Public = true // FIXME
	if !ValidateRoom(r) {
		return Room{}, ValidationErr
	}
	s.DB.Rooms = append(s.DB.Rooms, r)
	return r, nil
}

func (s *MonoService) GetRooms(page int, key string) ([]Room, int, error) {
	if key != "" && len(key) < 3 {
		return []Room{}, 0, TooShortKeyErr
	}
	if page <= 0 {
		return []Room{}, 0, OutOfRangeErr
	}
	if key == "" {
		firstRoomOnPageIdx := (page - 1) * 8
		lastRoomOnPageIdx := (page-1)*8 + 7
		if firstRoomOnPageIdx >= len(s.DB.Rooms)-1 && firstRoomOnPageIdx != 0 {
			return []Room{}, 0, OutOfRangeErr
		}
		return s.DB.Rooms[firstRoomOnPageIdx:min(len(s.DB.Rooms), lastRoomOnPageIdx+1)], (len(s.DB.Rooms)-1)/8 + 1, nil
	}
	filteredRooms := make([]Room, 0)
	for _, room := range s.DB.Rooms {
		if strings.Contains(strconv.FormatInt(room.ID, 10), key) || strings.Contains(room.PackageName, key) {
			filteredRooms = append(filteredRooms, room)
		}
	}
	firstRoomOnPageIdx := (page - 1) * 8
	lastRoomOnPageIdx := (page-1)*8 + 7
	if firstRoomOnPageIdx >= len(filteredRooms)-1 && firstRoomOnPageIdx != 0 {
		return []Room{}, 0, OutOfRangeErr
	}
	return filteredRooms[firstRoomOnPageIdx:min(len(filteredRooms), lastRoomOnPageIdx+1)], (len(filteredRooms)-1)/8 + 1, nil
}
