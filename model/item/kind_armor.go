package item

const (
	// KindArmor represents the kind of Armor
	KindArmor Kind = "armor"
)

// Armor describes the entity of an armor item
type Armor struct {
	Item `bson:",inline"`

	Type           string     `json:"type" bson:"type"`
	Armor          ArmorProps `json:"armor" bson:"armor"`
	RicochetChance string     `json:"ricochetChance,omitempty" bson:"ricochetChance,omitempty"`
	Penalties      Penalties  `json:"penalties" bson:"penalties"`
	Blocking       []string   `json:"blocking" bson:"blocking"`
	Slots          Slots      `json:"slots" bson:"slots"`
	Compatibility  List       `json:"compatibility" bson:"compatibility"`
}

// ArmorProps represents the armor properties of Armor and TacticalRig
type ArmorProps struct {
	Class           int64         `json:"class" bson:"class"`
	Durability      float64       `json:"durability" bson:"durability"`
	Material        ArmorMaterial `json:"material" bson:"material"`
	BluntThroughput float64       `json:"bluntThroughput" bson:"bluntThroughput"`
	Zones           []string      `json:"zones" bson:"zones"`
}

// ArmorMaterial represents the armor material of ArmorProps
type ArmorMaterial struct {
	Name            string  `json:"name" bson:"name"`
	Destructibility float64 `json:"destructibility" bson:"destructibility"`
}
