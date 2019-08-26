package model

import (
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ObjectID defines an object ID
type ObjectID = primitive.ObjectID

// Timestamp outputs the time in a Unix timestamp
type Timestamp struct {
	time.Time
}

// MarshalJSON implements the JSON marshaler
func (t *Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Unix())
}

// UnmarshalJSON implements the JSON unmarshaler
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

// NewResponse creates a new status response based on parameters
func NewResponse(msg string, code int) *Response {
	return &Response{
		Status:     http.StatusText(code),
		Message:    msg,
		StatusCode: code,
	}
}
