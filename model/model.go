package model

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
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

	if err := json.Unmarshal(b, &i); err != nil {
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

// Filter represents an MongoDB query filter
type Filter map[string]interface{}

// AddString adds a string to the given MongoDB field
func (f Filter) AddString(field, value string) error {
	var err error
	if value != "" {
		f[field], err = url.QueryUnescape(value)
	}

	return err
}

// AddInt adds an integer to the given MongoDB field
func (f Filter) AddInt(field, value string) error {
	if value != "" {
		var err error
		value, err = url.QueryUnescape(value)
		if err != nil {
			return err
		}

		f[field], err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddFloat adds a float to the given MongoDB field
func (f Filter) AddFloat(field, value string) error {
	if value != "" {
		var err error
		value, err = url.QueryUnescape(value)
		if err != nil {
			return err
		}

		f[field], err = strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
	}

	return nil
}
