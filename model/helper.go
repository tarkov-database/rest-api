package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ToObjectID converts a string id to an object ID
func ToObjectID(id string) (ObjectID, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return objID, ErrInvalidObjectID
	}

	return objID, nil
}
