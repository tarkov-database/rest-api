package item

import "go.mongodb.org/mongo-driver/bson"

const (
	// KindArmor represents the kind of Armor
	KindArmor Kind = "armor"
)

// Armor describes the entity of an armor item
type Armor struct {
	Item `bson:",inline"`

	Type           string           `json:"type" bson:"type"`
	Armor          ArmorProps       `json:"armor" bson:"armor"`
	Components     []ArmorComponent `json:"components" bson:"components"`
	RicochetChance string           `json:"ricochetChance,omitempty" bson:"ricochetChance,omitempty"`
	Penalties      Penalties        `json:"penalties" bson:"penalties"`
	Blocking       []string         `json:"blocking" bson:"blocking"`
	Slots          Slots            `json:"slots" bson:"slots"`
	Compatibility  List             `json:"compatibility" bson:"compatibility"`
	Conflicts      List             `json:"conflicts" bson:"conflicts"`
}

// ArmorComponent describes the entity of an armor component
type ArmorComponent struct {
	ArmorProps `bson:",inline"`
}

// ArmorProps represents the armor properties of ArmorComponent and Armor
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

// ArmorFilter describes the filters used for filtering Armor
type ArmorFilter struct {
	Type         *string
	ArmorClass   *int64
	MaterialName *string
}

// Filter implements the DocumentFilter interface
func (f *ArmorFilter) Filter() bson.D {
	filters := []bson.M{}

	if f.Type != nil {
		filters = append(filters, bson.M{"type": *f.Type})
	}

	if f.ArmorClass != nil {
		filters = append(filters, bson.M{"$or": []bson.M{
			{"armor.class": *f.ArmorClass},
			{"components": bson.M{"$elemMatch": bson.M{"class": *f.ArmorClass}}},
		}})
	}

	if f.MaterialName != nil {
		filters = append(filters, bson.M{"$or": []bson.M{
			{"armor.material.name": *f.MaterialName},
			{"components": bson.M{"$elemMatch": bson.M{"material.name": *f.MaterialName}}},
		}})
	}

	if len(filters) == 0 {
		return bson.D{}
	}

	return bson.D{{Key: "$and", Value: filters}}
}
