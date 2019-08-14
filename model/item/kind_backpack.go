package item

const (
	// KindBackpack represents the kind of Backpack
	KindBackpack Kind = "backpack"
)

// Backpack describes the entity of an backpack item
type Backpack struct {
	Item `json:",inline" bson:",inline"`

	Grids     []Grid    `json:"grids" bson:"grids"`
	Penalties Penalties `json:"penalties" bson:"penalties"`
}
