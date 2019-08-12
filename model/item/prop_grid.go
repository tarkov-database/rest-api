package item

type Grid struct {
	ID        string   `json:"id" bson:"id"`
	Height    int64    `json:"height" bson:"height"`
	Width     int64    `json:"width" bson:"width"`
	MaxWeight float64  `json:"maxWeight" bson:"maxWeight"`
	Filter    ItemList `json:"filter" bson:"filter"`
}
