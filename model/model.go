package model

import (
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type objectID = primitive.ObjectID

type Timestamp struct {
	time.Time
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Unix())
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	var i int64

	err := json.Unmarshal(b, &i)
	if err != nil {
		return err
	}

	*t = Timestamp{time.Unix(i, 0)}

	return nil
}

type Result struct {
	Count int64         `json:"total"`
	Items []interface{} `json:"items"`
}

type Response struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	StatusCode int    `json:"code"`
}

func (r *Response) New(msg string, code int) {
	r.Status = http.StatusText(code)
	r.Message = msg
	r.StatusCode = code
}
