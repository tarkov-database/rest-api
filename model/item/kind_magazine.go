package item

const (
	KindMagazine Kind = "magazine"
)

type Magazine struct {
	Item `json:",inline" bson:",inline"`

	Capacity      int64            `json:"capacity" bson:"capacity"`
	Caliber       string           `json:"caliber" bson:"caliber"`
	Ergonomics    int64            `json:"ergonomics" bson:"ergonomics"`
	Modifier      MagazineModifier `json:"modifier" bson:"modifier"`
	GridModifier  GridModifier     `json:"gridModifier" bson:"gridModifier"`
	Compatibility ItemList         `json:"compatibility" bson:"compatibility"`
}

type MagazineModifier struct {
	CheckTime  float64 `json:"checkTime" bson:"checkTime"`
	LoadUnload float64 `json:"loadUnload" bson:"loadUnload"`
}
