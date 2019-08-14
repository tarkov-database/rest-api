package item

// Grid describes the properties of a grid
type Grid struct {
	ID        string  `json:"id" bson:"id"`
	Height    int64   `json:"height" bson:"height"`
	Width     int64   `json:"width" bson:"width"`
	MaxWeight float64 `json:"maxWeight" bson:"maxWeight"`
	Filter    List    `json:"filter" bson:"filter"`
}
