package models

import "sigo/internal/lib"

type Statement struct {
	chooserId     int64
	stage         string
	usedQuestions [][]bool

	roundIdx int
	question *lib.Question
	slideIdx int
	/*	slideIdx     int           `json:"slide_idx"`
		answerableID int64         `json:"answerable"`*/
}

func (s *Statement) SlideIdx() int {
	return s.slideIdx
}

func (s *Statement) SetSlideIdx(slideIdx int) {
	s.slideIdx = slideIdx
}

func (s *Statement) Question() *lib.Question {
	return s.question
}

func (s *Statement) SetQuestion(question *lib.Question) {
	s.question = question
}

func (s *Statement) RoundIdx() int {
	return s.roundIdx
}

func (s *Statement) SetRoundIdx(roundIdx int) {
	s.roundIdx = roundIdx
}

func (s *Statement) SetUsedQuestions(usedQuestions [][]bool) {
	s.usedQuestions = usedQuestions
}

func (s *Statement) SetUsedQuestion(themeIdx, questionIdx int) {
	s.usedQuestions[themeIdx][questionIdx] = true
}

func (s *Statement) UsedQuestions() [][]bool {
	return s.usedQuestions
}

func (s *Statement) SetChooserId(uid int64) {
	s.chooserId = uid
}

func (s *Statement) ChooserId() int64 {
	return s.chooserId
}

func (s *Statement) SetStage(stage string) {
	s.stage = stage
}

func (s *Statement) Stage() string {
	return s.stage
}
