package item

const (
	// KindContainer represents the kind of Container
	KindContainer Kind = "container"
)

// Container describes the entity of an container item
type Container struct {
	Item `json:",inline" bson:",inline"`

	Grids []Grid `json:"grids" bson:"grids"`
}
