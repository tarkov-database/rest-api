package item

const (
	// KindGrenade represents the kind of Grenade
	KindGrenade Kind = "grenade"
)

// Grenade describes the entity of an grenade item
type Grenade struct {
	Item `json:",inline" bson:",inline"`

	Type              string  `json:"type" bson:"type"`
	Delay             float64 `json:"delay" bson:"delay"`
	FragmentCount     float64 `json:"fragCount" bson:"fragCount"`
	MinDistance       float64 `json:"minDistance" bson:"minDistance"`
	MaxDistance       float64 `json:"maxDistance" bson:"maxDistance"`
	ContusionDistance float64 `json:"contusionDistance" bson:"contusionDistance"`
	Strength          float64 `json:"strength" bson:"strength"`
	EmitTime          float64 `json:"emitTime" bson:"emitTime"`
}
