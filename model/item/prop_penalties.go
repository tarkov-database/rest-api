package item

// Penalties holds all existing penaltie types
type Penalties struct {
	Mouse           float64 `json:"mouse,omitempty" bson:"mouse,omitempty"`
	Speed           float64 `json:"speed,omitempty" bson:"speed,omitempty"`
	ErgonomicsFloat float64 `json:"ergonomicsFP,omitempty" bson:"ergonomicsFP,omitempty"`
	Ergonomics      int64   `json:"ergonomics,omitempty" bson:"ergonomics,omitempty"` // Deprecated: replaced
	Deafness        string  `json:"deafness,omitempty" bson:"deafness,omitempty"`
}
