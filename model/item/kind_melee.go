package item

const (
	// KindMelee represents the kind of Melee
	KindMelee Kind = "melee"
)

// Melee describes the entity of an melee item
type Melee struct {
	Item `json:",inline" bson:",inline"`

	Slash MeleeAttack `json:"slash" bson:"slash"`
	Stab  MeleeAttack `json:"stab" bson:"stab"`
}

// MeleeAttack represents the slash and stab data of Melee
type MeleeAttack struct {
	Damage      float64 `json:"damage" bson:"damage"`
	Rate        float64 `json:"rate" bson:"rate"`
	Range       float64 `json:"range" bson:"range"`
	Consumption float64 `json:"consumption" bson:"consumption"`
}
