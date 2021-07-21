package item

const (
	// KindBackpack represents the kind of Backpack
	KindBackpack Kind = "backpack"
)

// Backpack describes the entity of an backpack item
type Backpack struct {
	Item `bson:",inline"`

	Capacity  int64     `json:"capacity" bson:"capacity"`
	Grids     []Grid    `json:"grids" bson:"grids"`
	Penalties Penalties `json:"penalties" bson:"penalties"`
}
