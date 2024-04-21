package api

type Request struct {
	Type string `json:"type"`
	Data Data   `json:"data,omitempty"`
}

type Response struct {
	Type string `json:"type"`
	Data Data   `json:"data,omitempty"`
}

type Data struct {
	ThemeIndex    uint `json:"themeIndex,omitempty"`
	QuestionIndex uint `json:"questionIndex,omitempty"`
	PlayerId      int  `json:"playerId,omitempty"`
	ScoreChanges  int  `json:"scoreChanges,omitempty"`
	ChooserID     int  `json:"chooserId,omitempty"`
}
