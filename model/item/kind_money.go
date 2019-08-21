package item

const (
	// KindMoney represents the kind of Money
	KindMoney Kind = "money"
)

// Money describes the entity of an money item
type Money struct {
	Item `bson:",inline"`
}
