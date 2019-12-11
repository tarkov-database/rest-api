package item

// Effects holds all effect types
type Effects struct {
	Energy            *Effect `json:"energy,omitempty" bson:"energy,omitempty"`
	Hydration         *Effect `json:"hydration,omitempty" bson:"hydration,omitempty"`
	Bloodloss         *Effect `json:"bloodloss,omitempty" bson:"bloodloss,omitempty"`
	Fracture          *Effect `json:"fracture,omitempty" bson:"fracture,omitempty"`
	Contusion         *Effect `json:"contusion,omitempty" bson:"contusion,omitempty"`
	Pain              *Effect `json:"pain,omitempty" bson:"pain,omitempty"`
	Toxication        *Effect `json:"toxication,omitempty" bson:"toxication,omitempty"`
	RadiationExposure *Effect `json:"radExposure,omitempty" bson:"radExposure,omitempty"`
	Mobility          *Effect `json:"mobility,omitempty" bson:"mobility,omitempty"`
	Recoil            *Effect `json:"recoil,omitempty" bson:"recoil,omitempty"`
	ReloadSpeed       *Effect `json:"reloadSpeed,omitempty" bson:"reloadSpeed,omitempty"`
	LootSpeed         *Effect `json:"lootSpeed,omitempty" bson:"lootSpeed,omitempty"`
	UnlockSpeed       *Effect `json:"unlockSpeed,omitempty" bson:"unlockSpeed,omitempty"`
	DestroyedPart     *Effect `json:"destroyedPart,omitempty" bson:"destroyedPart,omitempty"`
}

// Effect represents the properties of an effect
type Effect struct {
	ResourceCosts int64           `json:"resourceCosts" bson:"resourceCosts"`
	FadeIn        float64         `json:"fadeIn" bson:"fadeIn"`
	FadeOut       float64         `json:"fadeOut" bson:"fadeOut"`
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
