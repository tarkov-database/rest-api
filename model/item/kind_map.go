package item

const (
	KindMap Kind = "map"
)

type Map struct {
	Item `json:",inline" bson:",inline"`
}
