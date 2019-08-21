package item

const (
	// KindMap represents the kind of Map
	KindMap Kind = "map"
)

// Map describes the entity of an map item
type Map struct {
	Item `bson:",inline"`
}
