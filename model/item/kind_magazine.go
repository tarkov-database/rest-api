package item

const (
	// KindMagazine represents the kind of Magazine
	KindMagazine Kind = "magazine"
)

// Magazine describes the entity of an magazine item
type Magazine struct {
	Item `bson:",inline"`

	Capacity      int64            `json:"capacity" bson:"capacity"`
	Caliber       string           `json:"caliber" bson:"caliber"`
	Ergonomics    int64            `json:"ergonomics" bson:"ergonomics"`
	Modifier      MagazineModifier `json:"modifier" bson:"modifier"`
	GridModifier  GridModifier     `json:"gridModifier" bson:"gridModifier"`
	Compatibility List             `json:"compatibility" bson:"compatibility"`
}

// MagazineModifier describes the properties of Modifier in Magazine
type MagazineModifier struct {
	CheckTime  float64 `json:"checkTime" bson:"checkTime"`
	LoadUnload float64 `json:"loadUnload" bson:"loadUnload"`
}
