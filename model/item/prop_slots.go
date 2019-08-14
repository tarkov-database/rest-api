package item

// Slots stores the mod slots of an item
type Slots map[string]Slot

// Slot represents a mod slot of an item
type Slot struct {
	Filter   List `json:"filter" bson:"filter"`
	Required bool `json:"required" bson:"required"`
}
