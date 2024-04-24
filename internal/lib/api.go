/* test requests:

{
    "type": "chooseQuestion",
    "data": {
        "themeIndex": 1,
        "questionIndex": 0
    }
}

{
    "type": "pressButton",
}

*/

package lib

type Request struct {
	UID  int64  `json:"uid,omitempty"`
	Type string `json:"type"`
	Data Data   `json:"data,omitempty"`
}

type Response struct {
	UID  int64  `json:"uid,omitempty"`
	Type string `json:"type"`
	Data Data   `json:"data,omitempty"`
}

type Data struct {
	QuestionType  string `json:"questionType,omitempty"`
	Content       string `json:"content,omitempty"`
	Status        string `json:"status,omitempty"`
	ThemeIndex    uint   `json:"themeIndex"`
	QuestionIndex uint   `json:"questionIndex"`
	ScoreChanges  int    `json:"scoreChanges,omitempty"`
	ChooserID     string `json:"chooserId,omitempty"`
}
