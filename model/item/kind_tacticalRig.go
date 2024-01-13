package item

import "go.mongodb.org/mongo-driver/bson"

const (
	// KindTacticalrig represents the kind of TacticalRig
	KindTacticalrig Kind = "tacticalrig"
)

// TacticalRig describes the entity of an tactical rig item
type TacticalRig struct {
	Item `bson:",inline"`

	Capacity        int64            `json:"capacity" bson:"capacity"`
	Grids           []Grid           `json:"grids" bson:"grids"`
	Penalties       Penalties        `json:"penalties" bson:"penalties"`
	Armor           *ArmorProps      `json:"armor,omitempty" bson:"armor,omitempty"` // Deprecated
	ArmorComponents []ArmorComponent `json:"armorComponents,omitempty" bson:"armorComponents,omitempty"`
	IsPlateCarrier  bool             `json:"isPlateCarrier" bson:"isPlateCarrier"`
	Slots           Slots            `json:"slots" bson:"slots"`
}

// TacticalRigFilter describes the filters used for filtering TacticalRig
type TacticalRigFilter struct {
	IsPlateCarrier *bool
	IsArmored      *bool
	ArmorClass     *int64
	ArmorMaterial  *string
}

// Filter implements the DocumentFilter interface
func (f *TacticalRigFilter) Filter() bson.D {
	filters := []bson.M{}

	if f.IsPlateCarrier != nil {
		filters = append(filters, bson.M{"isPlateCarrier": *f.IsPlateCarrier})
	}

	if f.IsArmored != nil {
		filters = append(filters, bson.M{"armorComponents": bson.M{"$exists": *f.IsArmored}})
	}

	if f.ArmorClass != nil {
		filters = append(filters, bson.M{
			"armorComponents": bson.M{"$elemMatch": bson.M{"class": *f.ArmorClass}},
		})
	}

	if f.ArmorMaterial != nil {
		filters = append(filters, bson.M{
			"armorComponents": bson.M{"$elemMatch": bson.M{"material.name": *f.ArmorMaterial}},
		})
	}

	if len(filters) == 0 {
		return bson.D{}
	}

	return bson.D{{Key: "$and", Value: filters}}
}
