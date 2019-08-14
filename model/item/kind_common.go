package item

import (
	"errors"

	"github.com/tarkov-database/rest-api/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type objectID = primitive.ObjectID

type timestamp = model.Timestamp

// List holds entity IDs and the associated kind
type List map[Kind][]objectID

// KindCommon represents the kind of Item
const KindCommon Kind = "common"

// Item represents the basic data of item
type Item struct {
	ID          objectID  `json:"_id" bson:"_id"`
	Name        string    `json:"name" bson:"name"`
	ShortName   string    `json:"shortName" bson:"shortName"`
	Description string    `json:"description" bson:"description"`
	Price       int64     `json:"price" bson:"price"`
	Weight      float64   `json:"weight" bson:"weight"`
	MaxStack    int64     `json:"maxStack" bson:"maxStack"`
	Rarity      string    `json:"rarity" bson:"rarity"`
	Grid        GridProps `json:"grid" bson:"grid"`
	Modified    timestamp `json:"_modified" bson:"_modified"`
	Kind        Kind      `json:"_kind" bson:"_kind"`
}

// GetID returns the ID of the item
func (i *Item) GetID() objectID {
	return i.ID
}

// SetID sets an ID
func (i *Item) SetID(id objectID) {
	i.ID = id
}

// GetKind returns the Kind of the item
func (i *Item) GetKind() Kind {
	return i.Kind
}

// SetKind sets an kind
func (i *Item) SetKind(k Kind) {
	i.Kind = k
}

// GetModified returns the modified date of the item
func (i *Item) GetModified() timestamp {
	return i.Modified
}

// SetModified sets an modified date
func (i *Item) SetModified(t timestamp) {
	i.Modified = t
}

// Validate validates the item fields
func (i *Item) Validate() error {
	if len(i.Name) < 3 {
		return errors.New("name is too short or not set")
	}
	if len(i.ShortName) < 1 {
		return errors.New("short name is too short or not set")
	}
	if len(i.Description) < 8 {
		return errors.New("description is too short or not set")
	}
	if i.Price < 0 {
		return errors.New("no negative price allowed")
	}
	if i.Weight < 0 {
		return errors.New("weight is too low or not set")
	}
	if i.MaxStack < 1 {
		return errors.New("maximum stack is too low or not set")
	}
	if i.Rarity == "" {
		return errors.New("rarity is not set")
	}
	if !i.Kind.IsValid() {
		return model.ErrInvalidKind
	}

	return nil
}

// GridProps represents the grid properties of an item
type GridProps struct {
	Color  RGBA  `json:"color" bson:"color"`
	Height int64 `json:"height" bson:"height"`
	Width  int64 `json:"width" bson:"width"`
}

// RGBA represents a color in RGBA
type RGBA struct {
	R uint `json:"r" bson:"r"`
	G uint `json:"g" bson:"g"`
	B uint `json:"b" bson:"b"`
	A uint `json:"a" bson:"a"`
}
