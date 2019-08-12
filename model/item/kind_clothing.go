package item

const (
	KindClothing Kind = "clothing"
)

type Clothing struct {
	Item `json:",inline" bson:",inline"`

	Type      string    `json:"type" bson:"type"`
	Blocking  []string  `json:"blocking" bson:"blocking"`
	Penalties Penalties `json:"penalties" bson:"penalties"`
	Slots     Slots     `json:"slots" bson:"slots"`
}
