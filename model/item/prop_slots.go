package item

type Slots map[string]Slot

type Slot struct {
	Filter   ItemList `json:"filter" bson:"filter"`
	Required bool     `json:"required" bson:"required"`
}
