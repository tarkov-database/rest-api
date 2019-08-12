package item

type Penalties struct {
	Mouse      float64 `json:"mouse,omitempty" bson:"mouse,omitempty"`
	Speed      float64 `json:"speed,omitempty" bson:"speed,omitempty"`
	Ergonomics int64   `json:"ergonomics,omitempty" bson:"ergonomics,omitempty"`
	Deafness   string  `json:"deafness,omitempty" bson:"deafness,omitempty"`
}
