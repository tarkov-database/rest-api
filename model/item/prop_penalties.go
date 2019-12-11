package item

// Penalties holds all existing penaltie types
type Penalties struct {
	Mouse         float64 `json:"mouse,omitempty" bson:"mouse,omitempty"`
	Speed         float64 `json:"speed,omitempty" bson:"speed,omitempty"`
	Ergonomics    float64 `json:"ergonomicsFP,omitempty" bson:"ergonomicsFP,omitempty"`
	ErgonomicsOld int64   `json:"ergonomics,omitempty" bson:"ergonomics,omitempty"` // Deprecated
	Deafness      string  `json:"deafness,omitempty" bson:"deafness,omitempty"`
}
