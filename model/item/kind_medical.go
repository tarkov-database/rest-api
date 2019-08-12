package item

const (
	KindMedical Kind = "medical"
)

type Medical struct {
	Item `json:",inline" bson:",inline"`

	Type         string  `json:"type" bson:"type"`
	Resources    int64   `json:"resources" bson:"resources"`
	ResourceRate int64   `json:"resourceRate" bson:"resourceRate"`
	UseTime      float64 `json:"useTime" bson:"useTime"`
	Effects      Effects `json:"effects" bson:"effects"`
}
