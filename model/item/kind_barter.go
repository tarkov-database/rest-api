package item

const (
	// KindBarter represents the kind of Barter
	KindBarter Kind = "barter"
)

// Barter describes the entity of an barter item
type Barter struct {
	Item `json:",inline" bson:",inline"`
}
