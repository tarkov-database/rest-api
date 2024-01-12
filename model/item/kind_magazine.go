package item

import "go.mongodb.org/mongo-driver/bson"

const (
	// KindMagazine represents the kind of Magazine
	KindMagazine Kind = "magazine"
)

// Magazine describes the entity of an magazine item
type Magazine struct {
	Item `bson:",inline"`

	Capacity          int64            `json:"capacity" bson:"capacity"`
	Caliber           string           `json:"caliber" bson:"caliber"`
	ErgonomicsFloat   float64          `json:"ergonomicsFP" bson:"ergonomicsFP"`
	Ergonomics        int64            `json:"ergonomics" bson:"ergonomics"` // Deprecated: replaced
	MalfunctionChance float64          `json:"malfunctionChance" bson:"malfunctionChance"`
	Modifier          MagazineModifier `json:"modifier" bson:"modifier"`
	GridModifier      GridModifier     `json:"gridModifier" bson:"gridModifier"`
	Compatibility     List             `json:"compatibility" bson:"compatibility"`
	Conflicts         List             `json:"conflicts" bson:"conflicts"`
}

// MagazineModifier describes the properties of Modifier in Magazine
type MagazineModifier struct {
	CheckTime  float64 `json:"checkTime" bson:"checkTime"`
	LoadUnload float64 `json:"loadUnload" bson:"loadUnload"`
}

// MagazineFilter describes the filters used for filtering Magazine
type MagazineFilter struct {
	Caliber *string
}

// Filter implements the DocumentFilter interface
func (f *MagazineFilter) Filter() bson.D {
	filters := []bson.M{}

	if f.Caliber != nil {
		filters = append(filters, bson.M{"caliber": *f.Caliber})
	}

	if len(filters) == 0 {
		return bson.D{}
	}

	return bson.D{{Key: "$and", Value: filters}}
}
