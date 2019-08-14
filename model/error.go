package model

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// ErrNoResult indicates that no documents were found
	ErrNoResult = errors.New("no document found")

	// ErrInvalidInput indicates that an input was invalid
	ErrInvalidInput = errors.New("invalid input")

	// ErrInvalidKind indicates that a kind was invalid
	ErrInvalidKind = errors.New("invalid kind")

	// ErrInvalidObjectID indicates that an object ID was invalid
	ErrInvalidObjectID = errors.New("invalid resource id")

	// ErrInternalError indicates that there was an function or backend error
	ErrInternalError = errors.New("server or network error")
)

// MongoToAPIError converts an MongoDB error to an internal error
func MongoToAPIError(err error) error {
	switch err {
	case mongo.ErrNoDocuments:
		return ErrNoResult
	default:
		return ErrInternalError
	}
}
