package item

const (
	// KindFirearm represents the kind of Firearm
	KindFirearm Kind = "firearm"
)

// Firearm describes the entity of an firearm item
type Firearm struct {
	Item `json:",inline" bson:",inline"`

	Type              string   `json:"type" bson:"type"`
	Class             string   `json:"class" bson:"class"`
	Caliber           string   `json:"caliber" bson:"caliber"`
	RateOfFire        int64    `json:"rof" bson:"rof"`
	Action            string   `json:"action" bson:"action"`
	Modes             []string `json:"modes" bson:"modes"`
	Velocity          float64  `json:"velocity" bson:"velocity"`
	EffectiveDistance int64    `json:"effectiveDist" bson:"effectiveDist"`
	Ergonomics        int64    `json:"ergonomics" bson:"ergonomics"`
	FoldRectractable  bool     `json:"foldRectractable" bson:"foldRectractable"`
	RecoilVertical    int64    `json:"recoilVertical" bson:"recoilVertical"`
	RecoilHorizontal  int64    `json:"recoilHorizontal" bson:"recoilHorizontal"`
	Slots             Slots    `json:"slots" bson:"slots"`
}
