package item

const (
	KindTacticalrig Kind = "tacticalrig"
)

type TacticalRig struct {
	Item `json:",inline" bson:",inline"`

	Grids     []Grid      `json:"grids" bson:"grids"`
	Penalties Penalties   `json:"penalties" bson:"penalties"`
	Armor     *ArmorProps `json:"armor,omitempty" bson:"armor,omitempty"`
}
