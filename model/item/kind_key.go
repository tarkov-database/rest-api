package item

const (
	KindKey Kind = "key"
)

type Key struct {
	Item `json:",inline" bson:",inline"`

	Location string `json:"location" bson:"location"`
}
