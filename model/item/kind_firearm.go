package item

const (
	// KindFirearm represents the kind of Firearm
	KindFirearm Kind = "firearm"
)

// Firearm describes the entity of an firearm item
type Firearm struct {
	Item `bson:",inline"`

	Type               string   `json:"type" bson:"type"`
	Class              string   `json:"class" bson:"class"`
	Caliber            string   `json:"caliber" bson:"caliber"`
	RateOfFire         int64    `json:"rof" bson:"rof"`
	BurstRounds        int64    `json:"burstRounds,omitempty" bson:"burstRounds,omitempty"`
	Action             string   `json:"action" bson:"action"`
	Modes              []string `json:"modes" bson:"modes"`
	Velocity           float64  `json:"velocity" bson:"velocity"`
	EffectiveDistance  int64    `json:"effectiveDist" bson:"effectiveDist"`
	ErgonomicsFloat    float64  `json:"ergonomicsFP" bson:"ergonomicsFP"`
	Ergonomics         int64    `json:"ergonomics" bson:"ergonomics"` // Deprecated: replaced
	FoldRectractable   bool     `json:"foldRectractable" bson:"foldRectractable"`
	RecoilVertical     int64    `json:"recoilVertical" bson:"recoilVertical"`
	RecoilHorizontal   int64    `json:"recoilHorizontal" bson:"recoilHorizontal"`
	OperatingResources float64  `json:"operatingResources" bson:"operatingResources"`
	MalfunctionChance  float64  `json:"malfunctionChance" bson:"malfunctionChance"`
	DurabilityRatio    float64  `json:"durabilityRatio" bson:"durabilityRatio"`
	HeatFactor         float64  `json:"heatFactor" bson:"heatFactor"`
	HeatFactorByShot   float64  `json:"heatFactorByShot" bson:"heatFactorByShot"`
	CoolFactor         float64  `json:"coolFactor" bson:"coolFactor"`
	CoolFactorMods     float64  `json:"coolFactorMods" bson:"coolFactorMods"`
	CenterOfImpact     float64  `json:"centerOfImpact" bson:"centerOfImpact"`
	Slots              Slots    `json:"slots" bson:"slots"`
}
