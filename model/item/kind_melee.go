package item

const (
	KindMelee Kind = "melee"
)

type Melee struct {
	Item `json:",inline" bson:",inline"`

	Slash MeleeAttack `json:"slash" bson:"slash"`
	Stab  MeleeAttack `json:"stab" bson:"stab"`
}

type MeleeAttack struct {
	Damage      float64 `json:"damage" bson:"damage"`
	Rate        float64 `json:"rate" bson:"rate"`
	Range       float64 `json:"range" bson:"range"`
	Consumption float64 `json:"consumption" bson:"consumption"`
}
