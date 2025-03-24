/* test requests:

{
    "type": "selectQuestion",
    "data": {
        "themeIdx": 1,
        "questionIdx": 1
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
	Content       string   `json:"content,omitempty"`
	ContentType   string   `json:"content_type,omitempty"`
	UsedQuestions [][]bool `json:"used_questions,omitempty"`
	ThemeIdx      int      `json:"theme_idx"`
	QuestionIdx   int      `json:"question_idx"`
	Time          int64    `json:"time"`
}
