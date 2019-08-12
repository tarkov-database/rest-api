package item

const (
	KindBackpack Kind = "backpack"
)

type Backpack struct {
	Item `json:",inline" bson:",inline"`

	Grids     []Grid    `json:"grids" bson:"grids"`
	Penalties Penalties `json:"penalties" bson:"penalties"`
}
