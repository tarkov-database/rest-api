package item

// Effects holds all effect types
type Effects struct {
	Energy            *Effect  `json:"energy,omitempty" bson:"energy,omitempty"`
	EnergyRate        *Effect  `json:"energyRate,omitempty" bson:"energyRate,omitempty"`
	Hydration         *Effect  `json:"hydration,omitempty" bson:"hydration,omitempty"`
	HydrationRate     *Effect  `json:"hydrationRate,omitempty" bson:"hydrationRate,omitempty"`
	Stamina           *Effect  `json:"stamina,omitempty" bson:"stamina,omitempty"`
	StaminaRate       *Effect  `json:"staminaRate,omitempty" bson:"staminaRate,omitempty"`
	Health            *Effect  `json:"health,omitempty" bson:"health,omitempty"`
	HealthRate        *Effect  `json:"healthRate,omitempty" bson:"healthRate,omitempty"`
	Bloodloss         *Effect  `json:"bloodloss,omitempty" bson:"bloodloss,omitempty"`
	Fracture          *Effect  `json:"fracture,omitempty" bson:"fracture,omitempty"`
	Contusion         *Effect  `json:"contusion,omitempty" bson:"contusion,omitempty"`
	Pain              *Effect  `json:"pain,omitempty" bson:"pain,omitempty"`
	TunnelVision      *Effect  `json:"tunnelVision,omitempty" bson:"tunnelVision,omitempty"`
	Tremor            *Effect  `json:"tremor,omitempty" bson:"tremor,omitempty"`
	Toxication        *Effect  `json:"toxication,omitempty" bson:"toxication,omitempty"`
	RadiationExposure *Effect  `json:"radExposure,omitempty" bson:"radExposure,omitempty"`
	Mobility          *Effect  `json:"mobility,omitempty" bson:"mobility,omitempty"`
	Recoil            *Effect  `json:"recoil,omitempty" bson:"recoil,omitempty"`
	ReloadSpeed       *Effect  `json:"reloadSpeed,omitempty" bson:"reloadSpeed,omitempty"`
	LootSpeed         *Effect  `json:"lootSpeed,omitempty" bson:"lootSpeed,omitempty"`
	UnlockSpeed       *Effect  `json:"unlockSpeed,omitempty" bson:"unlockSpeed,omitempty"`
	DestroyedPart     *Effect  `json:"destroyedPart,omitempty" bson:"destroyedPart,omitempty"`
	WeightLimit       *Effect  `json:"weightLimit,omitempty" bson:"weightLimit,omitempty"`
	DamageModifier    *Effect  `json:"damageModifier,omitempty" bson:"damageModifier,omitempty"`
	Skill             []Effect `json:"skill,omitempty" bson:"skill,omitempty"`
}

// Effect represents the properties of an effect
type Effect struct {
	Name          string          `json:"name,omitempty" bson:"name,omitempty"`
	ResourceCosts int64           `json:"resourceCosts" bson:"resourceCosts"`
	FadeIn        float64         `json:"fadeIn" bson:"fadeIn"`
	FadeOut       float64         `json:"fadeOut" bson:"fadeOut"`
	Chance        float64         `json:"chance" bson:"chance"`
	Delay         float64         `json:"delay" bson:"delay"`
	Duration      float64         `json:"duration" bson:"duration"`
	Value         float64         `json:"value" bson:"value"`
	IsPercent     bool            `json:"isPercent" bson:"isPercent"`
	Removes       bool            `json:"removes" bson:"removes"`
	Penalties     EffectPenalties `json:"penalties" bson:"penalties"`
}

// EffectPenalties holds the effect penalties
type EffectPenalties struct {
	HealthMin float64 `json:"healthMin,omitempty" bson:"healthMin,omitempty"`
	HealthMax float64 `json:"healthMax,omitempty" bson:"healthMax,omitempty"`
}
