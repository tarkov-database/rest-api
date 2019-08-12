package item

const (
	KindArmor Kind = "armor"
)

type Armor struct {
	Item `json:",inline" bson:",inline"`

	Type          string     `json:"type" bson:"type"`
	Armor         ArmorProps `json:"armor" bson:"armor"`
	Penalties     Penalties  `json:"penalties" bson:"penalties"`
	Blocking      []string   `json:"blocking" bson:"blocking"`
	Slots         Slots      `json:"slots" bson:"slots"`
	Compatibility ItemList   `json:"compatibility" bson:"compatibility"`
}

type ArmorProps struct {
	Class           int64         `json:"class" bson:"class"`
	Durability      float64       `json:"durability" bson:"durability"`
	Material        ArmorMaterial `json:"material" bson:"material"`
	BluntThroughput float64       `json:"bluntThroughput" bson:"bluntThroughput"`
	Zones           []string      `json:"zones" bson:"zones"`
}

type ArmorMaterial struct {
	Name            string  `json:"name" bson:"name"`
	Destructibility float64 `json:"destructibility" bson:"destructibility"`
}
