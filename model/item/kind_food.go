package item

const (
	// KindFood represents the kind of Food
	KindFood Kind = "food"
)

// Food describes the entity of an food item
type Food struct {
	Item `json:",inline" bson:",inline"`

	Type      string  `json:"type" bson:"type"`
	Resources int64   `json:"resources" bson:"resources"`
	UseTime   float64 `json:"useTime" bson:"useTime"`
	Effects   Effects `json:"effects" bson:"effects"`
}
