package item

const (
	// KindAmmunition represents the kind of Ammunition
	KindAmmunition Kind = "ammunition"
)

// Ammunition describes the entity of an ammunition item
type Ammunition struct {
	Item `bson:",inline"`

	Caliber             string   `json:"caliber" bson:"caliber"`
	Type                string   `json:"type" bson:"type"`
	Tracer              bool     `json:"tracer" bson:"tracer"`
	TracerColor         string   `json:"tracerColor" bson:"tracerColor"`
	Subsonic            bool     `json:"subsonic" bson:"subsonic"`
	Velocity            float64  `json:"velocity" bson:"velocity"`
	BallisticCoeficient float64  `json:"ballisticCoef" bson:"ballisticCoef"`
	Damage              float64  `json:"damage" bson:"damage"`
	Penetration         float64  `json:"penetration" bson:"penetration"`
	ArmorDamage         float64  `json:"armorDamage" bson:"armorDamage"`
	Fragmentation       AmmoFrag `json:"fragmentation" bson:"fragmentation"`
	Pellets             int64    `json:"pellets,omitempty" bson:"pellets,omitempty"`
}

// AmmoFrag represents the fragmentation data of Ammunition
type AmmoFrag struct {
	Chance float64 `json:"chance" bson:"chance"`
	Min    int64   `json:"min" bson:"min"`
	Max    int64   `json:"max" bson:"max"`
}
