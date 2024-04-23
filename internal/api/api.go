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

package api

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2/log"
)

type Request struct {
	Type string `json:"type"`
	Data Data   `json:"data,omitempty"`
}

type Response struct {
	Type string `json:"type"`
	Data Data   `json:"data,omitempty"`
}

type Data struct {
	Type          string `json:"type,omitempty"`
	Content       string `json:"content,omitempty"`
	Status        string `json:"status,omitempty"`
	ThemeIndex    uint   `json:"themeIndex"`
	QuestionIndex uint   `json:"questionIndex"`
	PlayerId      string `json:"playerId,omitempty"`
	ScoreChanges  int    `json:"scoreChanges,omitempty"`
	ChooserID     string `json:"chooserId,omitempty"`
}

// todo: error
func (r Request) Marshall() []byte {
	marshal, err := json.Marshal(r)
	if err != nil {
		log.Errorf("Error marshalling request: %v", err)
	}
	return marshal
}

func (r Response) Marshall() []byte {
	marshal, err := json.Marshal(r)
	if err != nil {
		log.Errorf("Error marshalling response: %v", err)
	}
	return marshal
}

func ReadResponse(b []byte) *Response {
	r := new(Response)
	err := json.Unmarshal(b, r)
	if err != nil {
		log.Errorf("Error unmarshalling response: %v", err)
	}
	return r
}
