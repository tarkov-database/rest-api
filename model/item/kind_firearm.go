package item

import "go.mongodb.org/mongo-driver/bson"

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
	Manufacturer       string   `json:"manufacturer" bson:"manufacturer"`
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

// FirearmFilter describes the filters used for filtering Firearm
type FirearmFilter struct {
	Manufacturer *string
	Type         *string
	Class        *string
	Caliber      *string
}

// Filter implements the DocumentFilter interface
func (f *FirearmFilter) Filter() bson.D {
	filters := []bson.M{}

	if f.Manufacturer != nil {
		filters = append(filters, bson.M{"manufacturer": *f.Manufacturer})
	}

	if f.Type != nil {
		filters = append(filters, bson.M{"type": *f.Type})
	}

	if f.Class != nil {
		filters = append(filters, bson.M{"class": *f.Class})
	}

	if f.Caliber != nil {
		filters = append(filters, bson.M{"caliber": *f.Caliber})
	}

	if len(filters) == 0 {
		return bson.D{}
	}

	return bson.D{{Key: "$and", Value: filters}}
}
