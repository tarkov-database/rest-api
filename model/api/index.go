package api

import (
// "github.com/tarkov-database/rest-api/model"
)

// type timestamp = model.Timestamp

const collection = "api"

type Index struct {
	Version string `json:"version" bson:"version"`
	// Modified map[string]timestamp `json:"modified" bson:"modified"`
}

func GetIndex() (*Index, error) {
	index := &Index{Version}

	// endpoints, err := getAllEndpoints()
	// if err != nil {
	// 	return index, err
	// }
	//
	// for _, ep := range endpoints {
	// 	index.Modified[ep.Name] = ep.Modified
	// }

	return index, nil
}
