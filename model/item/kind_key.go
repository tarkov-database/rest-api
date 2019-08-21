package item

const (
	// KindKey represents the kind of Key
	KindKey Kind = "key"
)

// Key describes the entity of an key item
type Key struct {
	Item `bson:",inline"`

	Location string `json:"location" bson:"location"`
}
