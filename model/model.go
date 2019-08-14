package model

import (
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type objectID = primitive.ObjectID

// Timestamp outputs the time in a Unix timestamp
type Timestamp struct {
	time.Time
}

// MarshalJSON implements the JSON marshaller
func (t *Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Unix())
}

// UnmarshalJSON implements the JSON unmarshaller
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	var i int64

	err := json.Unmarshal(b, &i)
	if err != nil {
		return err
	}

	*t = Timestamp{time.Unix(i, 0)}

	return nil
}

// Result describes an result output
type Result struct {
	Count int64         `json:"total"`
	Items []interface{} `json:"items"`
}

// Response describes a status response
type Response struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	StatusCode int    `json:"code"`
}

// New fills the Response with given values
func (r *Response) New(msg string, code int) {
	r.Status = http.StatusText(code)
	r.Message = msg
	r.StatusCode = code
}
