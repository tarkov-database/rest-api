package item

const (
	// KindContainer represents the kind of Container
	KindContainer Kind = "container"
)

// Container describes the entity of an container item
type Container struct {
	Item `bson:",inline"`

	Capacity int64  `json:"capacity" bson:"capacity"`
	Grids    []Grid `json:"grids" bson:"grids"`
}
