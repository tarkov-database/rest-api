package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToObjectID(id string) (objectID, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return objID, ErrInvalidObjectID
	}

	return objID, nil
}
