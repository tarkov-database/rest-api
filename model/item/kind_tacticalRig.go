package item

const (
	// KindTacticalrig represents the kind of TacticalRig
	KindTacticalrig Kind = "tacticalrig"
)

// TacticalRig describes the entity of an tactical rig item
type TacticalRig struct {
	Item `bson:",inline"`

	Grids     []Grid      `json:"grids" bson:"grids"`
	Penalties Penalties   `json:"penalties" bson:"penalties"`
	Armor     *ArmorProps `json:"armor,omitempty" bson:"armor,omitempty"`
}
