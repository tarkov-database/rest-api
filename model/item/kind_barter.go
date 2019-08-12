package item

const (
	KindBarter Kind = "barter"
)

type Barter struct {
	Item `json:",inline" bson:",inline"`
}
