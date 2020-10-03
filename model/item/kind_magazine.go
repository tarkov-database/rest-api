package item

const (
	// KindMagazine represents the kind of Magazine
	KindMagazine Kind = "magazine"
)

// Magazine describes the entity of an magazine item
type Magazine struct {
	Item `bson:",inline"`

	Capacity        int64            `json:"capacity" bson:"capacity"`
	Caliber         string           `json:"caliber" bson:"caliber"`
	ErgonomicsFloat float64          `json:"ergonomicsFP" bson:"ergonomicsFP"`
	Ergonomics      int64            `json:"ergonomics" bson:"ergonomics"` // Deprecated
	Modifier        MagazineModifier `json:"modifier" bson:"modifier"`
	GridModifier    GridModifier     `json:"gridModifier" bson:"gridModifier"`
	Compatibility   List             `json:"compatibility" bson:"compatibility"`
	Conflicts       List             `json:"conflicts" bson:"conflicts"`
}

// MagazineModifier describes the properties of Modifier in Magazine
type MagazineModifier struct {
	CheckTime  float64 `json:"checkTime" bson:"checkTime"`
	LoadUnload float64 `json:"loadUnload" bson:"loadUnload"`
}
