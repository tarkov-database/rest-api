package model

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNoResult        = errors.New("no document found")
	ErrInvalidInput    = errors.New("invalid input")
	ErrInvalidKind     = errors.New("invalid kind")
	ErrInvalidObjectID = errors.New("invalid resource id")
	ErrInternalError   = errors.New("server or network error")
)

func MongoToAPIError(err error) error {
	switch err {
	case mongo.ErrNoDocuments:
		return ErrNoResult
	default:
		return ErrInternalError
	}
}
