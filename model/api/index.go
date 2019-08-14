package api

const collection = "api"

// Index represents the object of the API root endpoint
type Index struct {
	Version string `json:"version" bson:"version"`
}

// GetIndex returns the data of the API root endpoint
func GetIndex() (*Index, error) {
	index := &Index{Version}

	return index, nil
}
