package item

const (
	// KindClothing represents the kind of Clothing
	KindClothing Kind = "clothing"
)

// Clothing describes the entity of an clothing item
type Clothing struct {
	Item `json:",inline" bson:",inline"`

	Type      string    `json:"type" bson:"type"`
	Blocking  []string  `json:"blocking" bson:"blocking"`
	Penalties Penalties `json:"penalties" bson:"penalties"`
	Slots     Slots     `json:"slots" bson:"slots"`
}
