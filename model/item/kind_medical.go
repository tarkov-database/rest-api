package item

const (
	// KindMedical represents the kind of Medical
	KindMedical Kind = "medical"
)

// Medical describes the entity of an medical item
type Medical struct {
	Item `bson:",inline"`

	Type         string  `json:"type" bson:"type"`
	Resources    int64   `json:"resources" bson:"resources"`
	ResourceRate int64   `json:"resourceRate" bson:"resourceRate"`
	UseTime      float64 `json:"useTime" bson:"useTime"`
	Effects      Effects `json:"effects" bson:"effects"`
}
