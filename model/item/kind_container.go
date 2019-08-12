package item

const (
	KindContainer Kind = "container"
)

type Container struct {
	Item `json:",inline" bson:",inline"`

	Grids []Grid `json:"grids" bson:"grids"`
}
