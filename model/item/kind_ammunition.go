package item

import "go.mongodb.org/mongo-driver/bson"

const (
	// KindAmmunition represents the kind of Ammunition
	KindAmmunition Kind = "ammunition"
)

// Ammunition describes the entity of an ammunition item
type Ammunition struct {
	Item `bson:",inline"`

	Caliber             string                 `json:"caliber" bson:"caliber"`
	Type                string                 `json:"type" bson:"type"`
	Tracer              bool                   `json:"tracer" bson:"tracer"`
	TracerColor         string                 `json:"tracerColor" bson:"tracerColor"`
	Subsonic            bool                   `json:"subsonic" bson:"subsonic"`
	CasingMass          float64                `json:"casingMass" bson:"casingMass"`
	BulletMass          float64                `json:"bulletMass" bson:"bulletMass"`
	BulletDiameter      float64                `json:"bulletDiameter" bson:"bulletDiameter"`
	Velocity            float64                `json:"velocity" bson:"velocity"`
	BallisticCoeficient float64                `json:"ballisticCoef" bson:"ballisticCoef"`
	Retardation         float64                `json:"retardation" bson:"retardation"`
	Damage              float64                `json:"damage" bson:"damage"`
	Penetration         float64                `json:"penetration" bson:"penetration"`
	ArmorDamage         float64                `json:"armorDamage" bson:"armorDamage"`
	Fragmentation       AmmoFrag               `json:"fragmentation" bson:"fragmentation"`
	Effects             AmmoEffects            `json:"effects" bson:"effects"`
	Projectiles         int64                  `json:"projectiles" bson:"projectiles"`
	Pellets             int64                  `json:"pellets,omitempty" bson:"pellets,omitempty"` // Deprecated: no longer used
	MisfireChance       float64                `json:"misfireChance" bson:"misfireChance"`
	FailureToFeedChance float64                `json:"failureToFeedChance" bson:"failureToFeedChance"`
	WeaponModifier      WeaponModifier         `json:"weaponModifier" bson:"weaponModifier"`
	GrenadeProperties   *AmmoGrenadeProperties `json:"grenadeProps,omitempty" bson:"grenadeProps,omitempty"`
}

// AmmoFrag represents the fragmentation data of Ammunition
type AmmoFrag struct {
	Chance float64 `json:"chance" bson:"chance"`
	Min    int64   `json:"min" bson:"min"`
	Max    int64   `json:"max" bson:"max"`
}

// AmmoEffects holds the effects of Ammunition
type AmmoEffects struct {
	LightBleedingChance float64 `json:"lightBleedingChance,omitempty" bson:"lightBleedingChance,omitempty"`
	HeavyBleedingChance float64 `json:"heavyBleedingChance,omitempty" bson:"heavyBleedingChance,omitempty"`
}

// WeaponModifier contains the weapon modifiers of Ammunition
type WeaponModifier struct {
	Accuracy          float64 `json:"accuracy" bson:"accuracy"`
	Recoil            float64 `json:"recoil" bson:"recoil"`
	MalfunctionChance float64 `json:"malfunctionChance" bson:"malfunctionChance"` // Deprecated: replaced
	DurabilityBurn    float64 `json:"durabilityBurn" bson:"durabilityBurn"`
	HeatFactor        float64 `json:"heatFactor" bson:"heatFactor"`
}

// AmmoGrenadeProperties represents the grenade properties of Ammunition
type AmmoGrenadeProperties struct {
	Delay         float64 `json:"delay" bson:"delay"`
	FragmentCount float64 `json:"fragCount" bson:"fragCount"`
	MinRadius     float64 `json:"minRadius" bson:"minRadius"`
	MaxRadius     float64 `json:"maxRadius" bson:"maxRadius"`
}

// AmmunitionFilter describes the filters used for filtering Ammunition
type AmmunitionFilter struct {
	Caliber *string
	Type    *string
}

// Filter implements the DocumentFilter interface
func (f *AmmunitionFilter) Filter() bson.D {
	filters := []bson.M{}

	if f.Caliber != nil {
		filters = append(filters, bson.M{"caliber": *f.Caliber})
	}

	if f.Type != nil {
		filters = append(filters, bson.M{"type": *f.Type})
	}

	if len(filters) == 0 {
		return bson.D{}
	}

	return bson.D{{Key: "$and", Value: filters}}
}
