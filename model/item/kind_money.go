package item

const (
	KindMoney Kind = "money"
)

type Money struct {
	Item `json:",inline" bson:",inline"`
}
